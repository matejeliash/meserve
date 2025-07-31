package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func GetDiskStatus() (string, error) {
	path, err := os.Getwd()

	if err != nil {
		return "", err
	}

	var stat unix.Statfs_t

	unix.Statfs(path, &stat)

	//freeGb := float64(stat.Bavail*uint64(stat.Bsize)) / 1_000_000_000.0
	//totalGb := float64(stat.Blocks*uint64(stat.Bsize)) / 1_000_000_000.0
	freeGb := formatSize(int64(stat.Bavail) * stat.Bsize)
	totalGb := formatSize(int64(stat.Blocks) * stat.Bsize)

	return fmt.Sprintf("%s / %s ", freeGb, totalGb), nil

}
