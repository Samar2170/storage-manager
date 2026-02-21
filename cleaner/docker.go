package cleaner

import (
	"runtime"
	"storage-manager/utils"
)

func checkDockerData() (int64, error) {
	arch := runtime.GOARCH
	os := runtime.GOOS
	dockerDir := ""
	if os == "linux" && arch == "amd64" {
		dockerDir = "/var/lib/docker"
	} else if os == "linux" && arch == "arm64" {
		dockerDir = "/var/lib/docker"
	} else if os == "windows" && arch == "amd64" {
		dockerDir = "C:\\ProgramData\\Docker"
	} else if os == "windows" && arch == "arm64" {
		dockerDir = "C:\\ProgramData\\Docker"
	} else if os == "darwin" && arch == "amd64" {
		dockerDir = "/var/lib/docker"
	} else if os == "darwin" && arch == "arm64" {
		dockerDir = "/var/lib/docker"
	}

	dockerSize, err := utils.GetRecursiveDirSize(dockerDir)
	if err != nil {
		return 0, err
	}
	return dockerSize, nil
}
