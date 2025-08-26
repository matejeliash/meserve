package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/matejeliash/meserve/handlers"
	"github.com/matejeliash/meserve/sysinfo"
)

func pathExists(path string) bool {
	var info os.FileInfo
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func main() {

	//servePath := "."

	portFlag := flag.Int("port", 8080, "enter port for server to use")
	serveDirFlag := flag.String("serveDir", ".", "root directory from which files are served")
	enableUploadFlag := flag.Bool("enableUpload", false, "usage enables uploading of files using form")

	flag.Usage = func() {
		fmt.Println("meserve - a simple fileserver written in Go that allows file browsing, downloading and uploading.")
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
	flag.Parse()

	portInt := *portFlag
	portStr := strconv.Itoa(portInt)
	selectedDir := *serveDirFlag
	enabledUpload := *enableUploadFlag

	if !pathExists(selectedDir) {
		//fmt.Println("wrong --serveDir , dir does not exist")
		log.Fatalf("--serveDir error, cannot access `%s`", selectedDir)
	}

	diskStatus, err := sysinfo.GetDiskStatus(selectedDir)
	if err != nil {
		log.Fatalf("failed to get disk space: %v\n", err)
	}
	fmt.Println(diskStatus)

	sysinfo.PrintAllAddresses(portInt)

	http.HandleFunc("GET /", handlers.FileHandler(selectedDir, enabledUpload))
	http.HandleFunc("POST /", handlers.UploadStreamHandler(selectedDir))

	fmt.Printf("Serving directory %s\n", selectedDir)
	err = http.ListenAndServe(":"+portStr, nil)

	if err != nil {
		log.Fatal(err)
	}
}
