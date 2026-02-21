package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("storage-manager", "simple storage tools")
	sizeInfoCmd := parser.NewCommand("size-info", "Get size info of a directory")
	folder := sizeInfoCmd.String("F", "folder", &argparse.Options{Required: true})
	checkcleanupCmd := parser.NewCommand("check-cleanup", "Clean up unnecessary files")
	cleanupFolder := checkcleanupCmd.String("F", "folder", &argparse.Options{Required: true})

	cleanupCmd := parser.NewCommand("cleanup", "Clean up unnecessary files")

	fmt.Println(*folder)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	switch {
	case sizeInfoCmd.Happened():
		breakdown, err := GetDirSizeBreakdown(*folder)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(breakdown)
	case checkcleanupCmd.Happened():
		err := CheckCleanup(*cleanupFolder)
		if err != nil {
			fmt.Println(err)
			return
		}
	case cleanupCmd.Happened():
		cleanup()
	default:
		fmt.Println(parser.Usage(nil))
	}
}

type DirSizeBreakdown struct {
	DirSizes  map[string]float64
	FileSizes map[string]float64
}

func (d DirSizeBreakdown) String() string {
	var out string
	out += "Directory Sizes:\n"
	out += "----------------\n"
	for name, size := range d.DirSizes {
		out += fmt.Sprintf("%-30s: %.2f Mb\n", name, size)
	}
	out += "\nFile Sizes:\n"
	out += "-----------\n"
	for name, size := range d.FileSizes {
		out += fmt.Sprintf("%-30s: %.2f Mb\n", name, size)
	}
	return out
}

func GetRecursiveDirSize(folder string) (int64, error) {
	var size int64
	err := filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func GetDirSizeBreakdown(folder string) (DirSizeBreakdown, error) {
	files, err := os.ReadDir(folder)

	if err != nil {

		return DirSizeBreakdown{}, err

	}

	breakdown := DirSizeBreakdown{
		DirSizes:  make(map[string]float64),
		FileSizes: make(map[string]float64),
	}

	for _, file := range files {

		if file.Name()[0] == '.' {

			continue

		}

		fullPath := filepath.Join(folder, file.Name())

		if file.IsDir() {

			size, err := GetRecursiveDirSize(fullPath)

			if err != nil {

				continue

			}

			breakdown.DirSizes[file.Name()] = float64(size) / 1024 / 1024

		} else {

			info, err := file.Info()
			if err != nil {
				continue
			}
			breakdown.FileSizes[file.Name()] = float64(info.Size()) / 1024 / 1024
		}
	}

	return breakdown, nil
}
