package main

import (
	"fmt"
	"os"

	"github.com/Samar2170/storage-manager/cleaner"
	"github.com/Samar2170/storage-manager/organizer"
	"github.com/Samar2170/storage-manager/utils"

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

	organizeCmd := parser.NewCommand("organize-photos", "Organize photos into folders based on location, camera, and orientation")
	organizeSource := organizeCmd.String("S", "source", &argparse.Options{Required: true, Help: "Source folder containing photos"})
	organizeDest := organizeCmd.String("D", "dest", &argparse.Options{Required: true, Help: "Destination folder to copy organized photos"})

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
	case organizeCmd.Happened():
		err := organizer.OrganizePhotos(*organizeSource, *organizeDest)
		if err != nil {
			fmt.Println("Error organizing photos:", err)
			return
		}
		fmt.Println("Successfully organized photos.")
	default:
		fmt.Println(parser.Usage(nil))
	}
}
