package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"storage-manager/cleaner"
	"storage-manager/utils"

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
		err := cleaner.CheckCleanup(*cleanupFolder)
		if err != nil {
			fmt.Println(err)
			return
		}
	case cleanupCmd.Happened():
		cleaner.Cleanup()
	default:
		fmt.Println(parser.Usage(nil))
	}
}

type SizeEntry struct {
	Name string
	Size float64
}

type DirSizeBreakdown struct {
	DirSizes  []SizeEntry
	FileSizes []SizeEntry
}

func (d DirSizeBreakdown) String() string {
	var out string
	out += "Directory Sizes:\n"
	out += "----------------\n"
	for _, entry := range d.DirSizes {
		out += fmt.Sprintf("%-30s: %.2f Mb\n", entry.Name, entry.Size)
	}
	out += "\nFile Sizes:\n"
	out += "-----------\n"
	for _, entry := range d.FileSizes {
		out += fmt.Sprintf("%-30s: %.2f Mb\n", entry.Name, entry.Size)
	}
	return out
}

func GetDirSizeBreakdown(folder string) (DirSizeBreakdown, error) {
	files, err := os.ReadDir(folder)

	if err != nil {

		return DirSizeBreakdown{}, err

	}

	breakdown := DirSizeBreakdown{}
	totalSize := 0.0

	for _, file := range files {
		if file.Name()[0] == '.' {
			continue
		}

		fullPath := filepath.Join(folder, file.Name())

		if file.IsDir() {
			size, err := utils.GetRecursiveDirSize(fullPath)
			if err != nil {
				continue
			}
			breakdown.DirSizes = append(breakdown.DirSizes, SizeEntry{
				Name: file.Name(),
				Size: float64(size) / 1024 / 1024,
			})
			totalSize += float64(size) / 1024 / 1024
		} else {
			info, err := file.Info()
			if err != nil {
				continue
			}
			breakdown.FileSizes = append(breakdown.FileSizes, SizeEntry{
				Name: file.Name(),
				Size: float64(info.Size()) / 1024 / 1024,
			})
			totalSize += float64(info.Size()) / 1024 / 1024
		}
	}

	sort.Slice(breakdown.DirSizes, func(i, j int) bool {
		return breakdown.DirSizes[i].Size > breakdown.DirSizes[j].Size
	})

	sort.Slice(breakdown.FileSizes, func(i, j int) bool {
		return breakdown.FileSizes[i].Size > breakdown.FileSizes[j].Size
	})

	breakdown.DirSizes = append(breakdown.DirSizes, SizeEntry{
		Name: "Total",
		Size: totalSize,
	})
	return breakdown, nil
}
