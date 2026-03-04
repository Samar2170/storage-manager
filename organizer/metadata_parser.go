package organizer

import (
	"context"
	"fmt"
	"image"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func ReadMetadata(filePath string) (ImageMetadata, error) {
	f, err := os.Open(filePath)
	var imd ImageMetadata = NewImageMetadata()
	if err != nil {
		return imd, err
	}
	defer f.Close()
	config, format, err := image.DecodeConfig(f)
	if err != nil {
		return imd, err
	}
	if config.Height >= config.Width {
		imd.Orientation = "vertical"
	} else {
		imd.Orientation = "landscape"
	}
	imd.Format = format

	exif.RegisterParsers(mknote.All...)

	f.Seek(0, 0)
	x, err := exif.Decode(f)
	if err != nil {
		return imd, err
	}

	makeTag, _ := x.Get(exif.Make)
	modelTag, _ := x.Get(exif.Model)

	if makeTag != nil || modelTag != nil {
		imd.Clicked = true
		imd.Camera = makeTag.String()
	}

	// Check for Location
	lat, long, err := x.LatLong()
	if err == nil {
		imd.Latitude = lat
		imd.Longitude = long
		imd.Location = fmt.Sprintf("%.4f_%.4f", lat, long)
	}

	// Walk all tags
	metadata := make(map[string]any)
	w := &metaWalker{data: metadata}
	err = x.Walk(w)
	if err != nil {
		return imd, err
	}
	imd.Tags = metadata

	creationTime, err := x.DateTime()
	if err == nil {
		imd.CreationTime = creationTime
	}
	return imd, nil
}

type metaWalker struct {
	data map[string]any
}

func (m *metaWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	m.data[string(name)] = tag
	return nil
}

func GetVideoMetadata(path string) (VideoMetadata, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var vmd VideoMetadata

	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return vmd, err
	}
	vmd.Tags = data.Format.TagList
	vmd.Duration = data.Format.Duration()
	vmd.Width = data.FirstVideoStream().Width
	vmd.Height = data.FirstVideoStream().Height
	vmd.Location, _ = data.Format.TagList.GetString("com.apple.quicktime.location.ISO6709")
	vmd.LivePhotoAuto, _ = data.Format.TagList.GetString("com.apple.quicktime.live-photo.auto")
	encoder, err := data.Format.TagList.GetString("encoder")
	make, _ := data.Format.TagList.GetString("com.apple.quicktime.make")
	if make == "" {
		vmd.Clicked = false
	}
	vmd.Camera = make
	if err == nil {
		vmd.Encoder = encoder
	}
	creationTime, err := data.Format.TagList.GetString("creation_time")
	if err == nil {
		vmd.CreationTime, err = time.Parse(time.RFC3339, creationTime)
	}

	if vmd.Height >= vmd.Width {
		vmd.Orientation = "vertical"
	} else {
		vmd.Orientation = "landscape"
	}

	return vmd, nil
}
