package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/matejeliash/meserve/files"
	"github.com/matejeliash/meserve/sysinfo"
	"github.com/matejeliash/meserve/tmpl"
)

func FileHandler(baseDir string, enabledUpload bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decodedPath, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		path := filepath.Join(baseDir, decodedPath)
		info, err := os.Stat(path)
		if err != nil {
			http.Error(w, "File with this url not found", http.StatusNotFound)
			return
		}

		if info.IsDir() {

			//start := time.Now()
			fileInfos, err := files.GetFileInfos(path)
			//fmt.Println(time.Since(start))
			if err != nil {
				http.Error(w, "Failed to get files", http.StatusInternalServerError)
				return
			}

			files.SortFileInfos(fileInfos)

			diskStatus, err := sysinfo.GetDiskStatus(path)
			if err != nil {
				http.Error(w, "Failed to get disk info", http.StatusInternalServerError)
				return
			}

			tmpl, err := tmpl.GetTemplate()

			if err != nil {
				http.Error(w, "Failed to parse template", http.StatusInternalServerError)
				return
			}
			data := struct {
				Files         []files.FileInfo
				Path          string
				DiskStatus    string
				EnabledUpload bool
			}{
				Files:         fileInfos,
				Path:          r.URL.Path,
				DiskStatus:    diskStatus,
				EnabledUpload: enabledUpload,
			}

			err = tmpl.Execute(w, data)

			if err != nil {
				http.Error(w, "Failed to  execute template", http.StatusInternalServerError)
				return
			}
		} else {

			// It's a file — serve it directly
			http.ServeFile(w, r, path)
		}
	}
}

func UploadStreamHandler(baseDir string) http.HandlerFunc {
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
