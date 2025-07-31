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

func GetFileInfos(files []os.DirEntry) []FileInfo {

	var fileInfos []FileInfo
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		fileInfo.ModTime()
		fileInfos = append(fileInfos, FileInfo{
			Name:             file.Name(),
			Size:             fileInfo.Size(),
			HumanSize:        padRight(formatSize(fileInfo.Size()), 10),
			ModTime:          fileInfo.ModTime(),
			FormattedModTime: formatTime(fileInfo.ModTime()),

			IsDir: file.IsDir(),
			//Path: filepath.ToSlash(filepath.Join(decodedPath, file.Name())) + func() string {
			//Path: url.PathEscape(filepath.ToSlash(filepath.Join(r.URL.Path, file.Name()))) + func() string {
			Path: url.PathEscape(file.Name()) + func() string {
				if file.IsDir() {
					return "/"
				}
				return ""
			}(),
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
