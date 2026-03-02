package organizer

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

var validImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}
var validVideoExtensions = map[string]struct{}{
	".mp4": {},
	".mov": {},
}

func ReadExif(filePath string) (*exif.Exif, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}
	return x, nil
}

type AnalyzePhotoFilterArgs struct {
	Location    bool
	Camera      bool
	Orientation bool
}

func AnalyzePhoto(filePath string, filterArgs AnalyzePhotoFilterArgs) (location string, camera string, orientation string, err error) {
	// Default values
	location = "unknown_location"
	camera = "downloaded"
	orientation = "unknown_orientation"

	f, err := os.Open(filePath)
	if err != nil {
		return location, camera, orientation, err
	}
	defer f.Close()

	// 1. Determine Orientation
	imd, err := ReadMetadata(filePath)
	if err != nil {
		return location, camera, orientation, err
	}
	if filterArgs.Location {
		location = imd.Location
	}
	if filterArgs.Camera {
		camera = imd.Camera
	}
	if filterArgs.Orientation {
		orientation = imd.Orientation
	}
	return location, camera, orientation, nil
}

func AnalyzeVideo(filePath string, filterArgs AnalyzePhotoFilterArgs) (location string, camera string, orientation string, err error) {
	// Default values
	location = "unknown_location"
	camera = "downloaded"
	orientation = "unknown_orientation"

	f, err := os.Open(filePath)
	if err != nil {
		return location, camera, orientation, err
	}
	defer f.Close()

	// 1. Determine Orientation
	vmd, err := GetVideoMetadata(filePath)
	if err != nil {
		return location, camera, orientation, err
	}
	if filterArgs.Location {
		location = vmd.Location
	}
	// if filterArgs.Camera {
	// camera = vmd.
	// }
	if filterArgs.Orientation {
		orientation = vmd.Orientation
	}
	return location, camera, orientation, nil
}

func OrganizePhotos(sourceDir string, destDir string, filterArgs AnalyzePhotoFilterArgs) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		validExtensions := validImageExtensions
		for k, v := range validVideoExtensions {
			validExtensions[k] = v
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := validExtensions[ext]; !ok {
			return nil
		}

		location, camera, orientation, err := AnalyzePhoto(path, filterArgs)
		if err != nil {
			log.Printf("Failed to analyze %s: %v\n", path, err)
			// Decide to skip or put in an "unknown" directory. Let's still try to organize it.
		}

		// Build destination path
		destPath := filepath.Join(destDir, location, camera, orientation)
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			log.Printf("Failed to create directory %s: %v\n", destPath, err)
			return err
		}

		destFile := filepath.Join(destPath, info.Name())
		err = copyFile(path, destFile)
		if err != nil {
			log.Printf("Failed to copy file %s to %s: %v\n", path, destFile, err)
			return err
		}

		fmt.Printf("Copied %s -> %s\n", path, destFile)
		return nil
	})
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
