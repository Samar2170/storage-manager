package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	removeFileTimeConst = 12 * time.Hour
)

var cleanupStuff = map[string]struct{}{
	"venv":         {},
	"venv*":        {},
	"*venv":        {},
	"*env":         {},
	"__pycache__":  {},
	".git":         {},
	".idea":        {},
	".vscode":      {},
	".DS_Store":    {},
	"node_modules": {},
	"dist":         {},
	"build":        {},
	"*.log":        {},
	"*.tmp":        {},
	"*.bak":        {},
	"*.swp":        {},
	"*.swo":        {},
	"*.pyc":        {},
	"*.pyo":        {},
	"*.pyd":        {},
}

var packageFiles = map[string]struct{}{
	"*.egg-info":  {},
	"*.dist-info": {},
	"*tar.gz":     {},
	"*whl":        {},
	"*.dmg":       {},
	"*.iso":       {},
	"*.exe":       {},
	"*.msi":       {},
	"*.deb":       {},
	"*.rpm":       {},
	"*.apk":       {},
	"*.ipa":       {},
}

func shouldCleanup(name string) bool {
	for pattern := range cleanupStuff {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	for pattern := range packageFiles {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

func findCleanupEntries(folder string) (float64, map[string]float64, error) {
	var entries = make(map[string]float64)
	var totalSize float64

	var dirQueue []string
	dirQueue = append(dirQueue, folder)
	for len(dirQueue) > 0 {
		currFolder := dirQueue[0]
		dirQueue = dirQueue[1:]
		files, err := os.ReadDir(currFolder)
		if err != nil {
			return totalSize, entries, err
		}
		for _, d := range files {
			fullPath := filepath.Join(currFolder, d.Name())
			if shouldCleanup(d.Name()) {
				var size int64
				if d.IsDir() {
					size, err = GetRecursiveDirSize(fullPath)
				} else {
					info, err := d.Info()
					if err != nil {
						return totalSize, entries, err
					}
					size = info.Size()
				}

				if err != nil {
					return totalSize, entries, err
				}

				sizeMb := float64(size) / 1024 / 1024
				entries[fullPath] = sizeMb
				totalSize += sizeMb
				continue
			}

			if d.IsDir() {
				if d.Name()[0] == '.' || shouldCleanup(d.Name()) {
					continue
				}
				dirQueue = append(dirQueue, fullPath)
			}
		}

	}
	return totalSize, entries, nil
}

func cleanupFileManager(entries map[string]float64, totalSize float64) {
	// check previous file, if more than a day old, remove it
	files, err := os.ReadDir(".")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), "cleanup_file_") && strings.HasSuffix(file.Name(), ".tmp") {
				info, err := file.Info()
				if err == nil {
					if time.Since(info.ModTime()) > removeFileTimeConst {
						os.Remove(file.Name())
					}
				}
			}
		}
	}
	fileData := map[string]interface{}{
		"totalSize": totalSize,
		"entries":   entries,
	}
	// save cleanup_file as tmp with date and timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("cleanup_file_%s.tmp", timestamp)
	data, err := json.MarshalIndent(fileData, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling entries: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing cleanup file: %v\n", err)
	} else {
		fmt.Printf("Cleanup entries saved to %s\n", filename)
	}
}

func CheckCleanup(folder string) error {
	totalSize, entries, err := findCleanupEntries(folder)

	if err != nil {
		return err
	}
	for path, size := range entries {
		fmt.Printf("Path: %s, Size: %.2f Mb\n", path, size)
	}
	fmt.Printf("Total Size: %.2f Mb\n", totalSize)

	cleanupFileManager(entries, totalSize)
	return nil
}

func cleanup() {
	files, err := os.ReadDir(".")
	if err != nil {
		return
	}
	cleanupFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "cleanup_file_") && strings.HasSuffix(file.Name(), ".tmp") {
			info, err := file.Info()
			if err == nil {
				if time.Since(info.ModTime()) > 24*time.Hour {
					os.Remove(file.Name())
				} else {
					cleanupFiles = append(cleanupFiles, file.Name())
				}
			}
		}
	}

	for c, file := range cleanupFiles {
		var fileData map[string]interface{}
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		err = json.Unmarshal(data, &fileData)
		if err != nil {
			continue
		}
		fmt.Printf("Cleanup File %d:\n", c+1)
		for path := range fileData["entries"].(map[string]interface{}) {
			fmt.Println(path)
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Printf("Error removing %s: %v\n", path, err)
			}
		}
		fmt.Printf("Cleaned up %.2f Mb\n", fileData["totalSize"])
		os.Remove(file)
	}
}
