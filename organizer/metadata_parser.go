package organizer

import (
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"github.com/rwcarlsen/goexif/tiff"
)

func ReadMetadata(filePath string) (map[string]any, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Walk all tags
	metadata := make(map[string]any)
	w := &metaWalker{data: metadata}
	err = x.Walk(w)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

type metaWalker struct {
	data map[string]any
}

func (m *metaWalker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	m.data[string(name)] = tag
	return nil
}
