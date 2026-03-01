package organizer

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

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

func AnalyzePhoto(filePath string) (location string, camera string, orientation string, err error) {
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
	config, format, err := image.DecodeConfig(f)
	if err == nil && (format == "jpeg" || format == "png") {
		if config.Height >= config.Width {
			orientation = "vertical"
		} else {
			orientation = "landscape"
		}
	}

	// 2. Determine Camera and Location
	// Rewind the file to the beginning so goexif can read it
	f.Seek(0, 0)
	x, err := exif.Decode(f)
	if err == nil {
		// Check for Make or Model tags to classify as "clicked"
		makeTag, _ := x.Get(exif.Make)
		modelTag, _ := x.Get(exif.Model)

		if makeTag != nil || modelTag != nil {
			camera = "clicked"
		}

		// Check for Location
		lat, long, err := x.LatLong()
		if err == nil {
			location = fmt.Sprintf("%.4f_%.4f", lat, long)
		}
	}

	return location, camera, orientation, nil
}

func OrganizePhotos(sourceDir string, destDir string) error {
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			// Skip non-image files
			return nil
		}

		location, camera, orientation, err := AnalyzePhoto(path)
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
