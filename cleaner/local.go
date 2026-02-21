package cleaner

import (
	"runtime"
	"storage-manager/utils"
)

func readLocalData() (int64, error) {
	arch := runtime.GOARCH
	os := runtime.GOOS
	localFolder := ""
	if arch == "amd64" && os == "linux" {
		localFolder = "~/.local/share"
	}
	return utils.GetRecursiveDirSize(localFolder)

}
