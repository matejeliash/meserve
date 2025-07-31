package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// func uploadHandler(baseDir string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		decodedPath, err := url.PathUnescape(r.URL.Path)
// 		fmt.Println(decodedPath)
// 		if err != nil {
// 			http.Error(w, "Invalid URL path.", http.StatusBadRequest)
// 			return
// 		}
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		err = r.ParseMultipartForm(64 << 20) // 32MB
// 		if err != nil {
// 			http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
// 			return
// 		}

// 		files := r.MultipartForm.File["file"]
// 		if len(files) == 0 {
// 			http.Error(w, "No files uploaded", http.StatusBadRequest)
// 			return
// 		}

// 		for _, fileHeader := range files {
// 			file, err := fileHeader.Open()
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Failed to open file: %s", fileHeader.Filename), http.StatusInternalServerError)
// 				return
// 			}
// 			//defer file.Close()

// 			path := filepath.Join(baseDir, decodedPath)
// 			dstPath := filepath.Join(path, fileHeader.Filename)
// 			fmt.Println(path)
// 			fmt.Println(dstPath)
// 			outFile, err := os.Create(dstPath)
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Failed to create file: %s", dstPath), http.StatusInternalServerError)
// 				return
// 			}
// 			//defer outFile.Close()

// 			_, err = io.Copy(outFile, file)
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("Failed to save file: %s", dstPath), http.StatusInternalServerError)
// 				return
// 			}

// 			file.Close()
// 			outFile.Close()

// 			log.Printf("Saved file: %s", dstPath)
// 		}

// 	}
// }

func main() {

	portPtr := flag.Int("port", 8080, "enter port for server to use")

	flag.Parse()
	//port := ":8080"
	//
	portInt := *portPtr
	portStr := strconv.Itoa(portInt)

	PrintAllAddresses(portInt)

	diskStatus, err := GetDiskStatus()
	if err != nil {
		fmt.Println("Failed to get disk status:", err)
	}
	fmt.Println(diskStatus)

	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}

	http.HandleFunc("GET /", fileHandler(baseDir))
	http.HandleFunc("POST /", uploadStreamHandler(baseDir))

	fmt.Printf("Serving directory %s\n", baseDir)
	err = http.ListenAndServe(":"+portStr, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
}
