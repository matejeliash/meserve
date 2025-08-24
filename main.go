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
	flag.Parse()
	portInt := *portFlag
	portStr := strconv.Itoa(portInt)

	// TODO add port range checking

	// baseDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Failed to get current directory:", err)
	// 	return
	// }

	selectedDir := *serveDirFlag
	if !pathExists(selectedDir) {
		//fmt.Println("wrong --serveDir , dir does not exist")
		log.Fatalf("--serveDir error, cannot access `%s`", selectedDir)

	}

	sysinfo.PrintAllAddresses(portInt)

	diskStatus, err := sysinfo.GetDiskStatus(selectedDir)
	if err != nil {
		log.Fatalf("failder to get disk space: %v\n", err)
	}
	fmt.Println(diskStatus)

	http.HandleFunc("GET /", handlers.FileHandler(selectedDir))
	http.HandleFunc("POST /", handlers.UploadStreamHandler(selectedDir))

	fmt.Printf("Serving directory %s\n", selectedDir)
	err = http.ListenAndServe(":"+portStr, nil)

	if err != nil {
		log.Fatal(err)
	}
}
