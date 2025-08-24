package sysinfo

import (
	"fmt"

	"github.com/matejeliash/meserve/format"
	"golang.org/x/sys/unix"
)

func GetDiskStatus(path string) (string, error) {

	var stat unix.Statfs_t

	unix.Statfs(path, &stat)

	//freeGb := float64(stat.Bavail*uint64(stat.Bsize)) / 1_000_000_000.0
	//totalGb := float64(stat.Blocks*uint64(stat.Bsize)) / 1_000_000_000.0
	freeGb := format.FormatSize(int64(stat.Bavail) * stat.Bsize)
	totalGb := format.FormatSize(int64(stat.Blocks) * stat.Bsize)

	return fmt.Sprintf("%s / %s ", freeGb, totalGb), nil

}
