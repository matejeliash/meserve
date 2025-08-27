package files

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/matejeliash/meserve/internal/format"
)

type FileInfo struct {
	Name             string
	HumanSize        string
	Size             int64
	IsDir            bool
	Path             string
	ModTime          time.Time
	UnixModTime      int64
	FormattedModTime string
}

func GetFileInfos(path string) ([]FileInfo, error) {

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]FileInfo, 0, len(files))

	for _, file := range files {
		// Get basic info (does not follow symlinks)
		lstatInfo, err := file.Info()
		if err != nil {
			continue
		}

		// Resolve symlinks using os.Stat (follows symlinks)
		fullPath := filepath.Join(path, file.Name())
		statInfo, err := os.Stat(fullPath)

		if err != nil {
			// fallback: use lstatInfo if stat fails
			statInfo = lstatInfo
			fmt.Println("revert to lstat")
		}

		isDir := statInfo.IsDir()

		escapedPath := url.PathEscape(file.Name())
		if isDir {
			escapedPath += "/"
		}

		fileInfos = append(fileInfos, FileInfo{
			Name:             file.Name(),
			Size:             statInfo.Size(),
			HumanSize:        format.PadRight(format.FormatSize(statInfo.Size()), 20),
			ModTime:          statInfo.ModTime(),
			FormattedModTime: format.FormatTime(statInfo.ModTime()),
			IsDir:            isDir,
			UnixModTime:      statInfo.ModTime().UnixMilli(),
			Path:             escapedPath,
		})
	}

	return fileInfos, nil
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
