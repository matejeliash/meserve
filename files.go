package main

import (
	"net/url"
	"os"
	"sort"
	"time"
)

type FileInfo struct {
	Name             string
	HumanSize        string
	Size             int64
	IsDir            bool
	Path             string
	ModTime          time.Time
	FormattedModTime string
}

// change this so it uses os.Stat
func GetFileInfos(files []os.DirEntry) []FileInfo {

	fileInfos := make([]FileInfo, 0, len(files))
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue

		}
		fileInfos = append(fileInfos, FileInfo{
			Name:             file.Name(),
			Size:             fileInfo.Size(),
			HumanSize:        padRight(formatSize(fileInfo.Size()), 10),
			ModTime:          fileInfo.ModTime(),
			FormattedModTime: formatTime(fileInfo.ModTime()),

			IsDir: fileInfo.IsDir(),

			Path: url.PathEscape(file.Name()) + "/",
		})
	}
	return fileInfos
}

func SortFileInfos(files []FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]

		// Directories first
		if a.IsDir && !b.IsDir {
			return true
		}
		if !a.IsDir && b.IsDir {
			return false
		}

		// Both directories: sort by name ascending
		if a.IsDir && b.IsDir {
			return a.Name < b.Name
		}

		// Both files: sort by size descending
		return a.Size > b.Size
	})
}
