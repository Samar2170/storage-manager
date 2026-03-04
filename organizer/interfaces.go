package organizer

import (
	"fmt"
	"strings"
	"time"
)

type ImageMetadata struct {
	Orientation  string
	Camera       string
	Clicked      bool
	Location     string
	Format       string
	Tags         map[string]any
	Latitude     float64
	Longitude    float64
	CreationTime time.Time
}

func NewImageMetadata() ImageMetadata {
	return ImageMetadata{
		Orientation:  "unknown_orientation",
		Camera:       "downloaded",
		Location:     "unknown_location",
		Format:       "unknown_format",
		Tags:         make(map[string]any),
		Latitude:     0,
		Longitude:    0,
		CreationTime: time.Time{},
	}
}

type VideoMetadata struct {
	Width         int
	Height        int
	Location      string
	LivePhotoAuto string
	Encoder       string
	Duration      time.Duration
	CreationTime  time.Time
	Tags          map[string]any
	Orientation   string
	Camera        string
	Clicked       bool
}

func NewVideoMetadata() VideoMetadata {
	return VideoMetadata{
		Orientation: "unknown_orientation",
		Camera:      "downloaded",
		Location:    "unknown_location",
		Tags:        make(map[string]any),
	}
}

func (vm VideoMetadata) String() string {
	var sb strings.Builder
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	sb.WriteString("| Field                | Value                                              |\n")
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	sb.WriteString(fmt.Sprintf("| Width                | %-50d |\n", vm.Width))
	sb.WriteString(fmt.Sprintf("| Height               | %-50d |\n", vm.Height))
	sb.WriteString(fmt.Sprintf("| Location             | %-50s |\n", vm.Location))
	sb.WriteString(fmt.Sprintf("| LivePhotoAuto        | %-50s |\n", vm.LivePhotoAuto))
	sb.WriteString(fmt.Sprintf("| Encoder              | %-50s |\n", vm.Encoder))
	sb.WriteString(fmt.Sprintf("| Duration             | %-50s |\n", vm.Duration))
	sb.WriteString(fmt.Sprintf("| CreationTime         | %-50s |\n", vm.CreationTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("| Orientation          | %-50s |\n", vm.Orientation))
	sb.WriteString(fmt.Sprintf("| Camera               | %-50s |\n", vm.Camera))
	sb.WriteString("+----------------------+----------------------------------------------------+\n")

	for key, value := range vm.Tags {
		sb.WriteString(fmt.Sprintf("| %-20s | %-50s |\n", key, value))
	}
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	return sb.String()
}

func (im ImageMetadata) String() string {
	var sb strings.Builder
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	sb.WriteString("| Field                | Value                                              |\n")
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	sb.WriteString(fmt.Sprintf("| Orientation          | %-50s |\n", im.Orientation))
	sb.WriteString(fmt.Sprintf("| Camera               | %-50s |\n", im.Camera))
	sb.WriteString(fmt.Sprintf("| Location             | %-50s |\n", im.Location))
	sb.WriteString(fmt.Sprintf("| Format               | %-50s |\n", im.Format))
	sb.WriteString(fmt.Sprintf("| Latitude             | %-50.6f |\n", im.Latitude))
	sb.WriteString(fmt.Sprintf("| Longitude            | %-50.6f |\n", im.Longitude))
	sb.WriteString(fmt.Sprintf("| CreationTime         | %-50s |\n", im.CreationTime.Format("2006-01-02 15:04:05")))
	sb.WriteString("+----------------------+----------------------------------------------------+\n")

	for key, value := range im.Tags {
		sb.WriteString(fmt.Sprintf("| %-20s | %-50s |\n", key, value))
	}
	sb.WriteString("+----------------------+----------------------------------------------------+\n")
	return sb.String()
}
