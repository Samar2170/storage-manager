package main

import (
	"fmt"
	"os"
	"storage-manager/cleaner"
	"storage-manager/organizer"
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

	metadataCmd := parser.NewCommand("metadata", "Read metadata of a file")
	metadataFile := metadataCmd.String("F", "file", &argparse.Options{Required: true})

	fmt.Println(*folder)
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	switch {
	case sizeInfoCmd.Happened():
		breakdown, err := utils.GetDirSizeBreakdown(*folder)
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

	case metadataCmd.Happened():
		_, err := organizer.ReadMetadata(*metadataFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println(parser.Usage(nil))
	}
}
