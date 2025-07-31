package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func fileHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decodedPath, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URL path.", http.StatusBadRequest)
			return
		}

		path := filepath.Join(baseDir, decodedPath)
		info, err := os.Stat(path)
		if err != nil {
			http.Error(w, "file with this url not found", http.StatusNotFound)
			return
		}

		if info.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				http.Error(w, "Cannot read directory.", http.StatusInternalServerError)
				return
			}

			fileInfos := GetFileInfos(files)

			SortFileInfos(fileInfos)

			diskStatus, err := GetDiskStatus()
			if err != nil {
				http.Error(w, "Cannot get diskstatus.", http.StatusInternalServerError)
				return
			}

			tmpl := GetTemplate()
			data := struct {
				Files      []FileInfo
				Path       string
				DiskStatus string
			}{
				Files:      fileInfos,
				Path:       r.URL.Path,
				DiskStatus: diskStatus,
			}

			tmpl.Execute(w, data)
		} else {
			// It's a file — serve it directly
			http.ServeFile(w, r, path)
		}
	}
}

func uploadStreamHandler(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decodedPath, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URL path.", http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mr, err := r.MultipartReader()
		if err != nil {
			http.Error(w, "Failed to create multipart reader", http.StatusBadRequest)
			return
		}

		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Failed to read part", http.StatusInternalServerError)
				return
			}
			if part.FileName() == "" {
				continue // Skip non-file fields
			}

			dstPath := filepath.Join(baseDir, decodedPath, filepath.Base(part.FileName()))
			outFile, err := os.Create(dstPath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create file: %s", dstPath), http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(outFile, part)
			outFile.Close()
			part.Close()

			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to save file: %s", dstPath), http.StatusInternalServerError)
				return
			}
			log.Printf("Saved file: %s", dstPath)
		}
	}
}

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
