package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/matejeliash/meserve/internal/files"
	"github.com/matejeliash/meserve/internal/sysinfo"
	"github.com/matejeliash/meserve/internal/tmpl"
)

func FileHandler(baseDir string, enabledUpload bool, enabledDiskStatus bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decodedPath, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}

		decodedPath = strings.TrimPrefix(decodedPath, "/")
		realPath := filepath.Join(baseDir, decodedPath)

		realPath = strings.ReplaceAll(realPath, "-_PERCENT_-", "%")

		// fmt.Println(path)
		info, err := os.Stat(realPath)
		if err != nil {
			http.Error(w, "File with this url not found", http.StatusNotFound)
			return
		}

		if info.IsDir() {

			// If file in URL is dir for add / the end  by redirection
			if !strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
				return
			}

			//start := time.Now()
			fileInfos, err := files.GetFileInfos(realPath)
			//fmt.Println(time.Since(start))
			if err != nil {
				http.Error(w, "Failed to get files", http.StatusInternalServerError)
				return
			}

			files.SortFileInfos(fileInfos)

			diskStatus := ""

			if enabledUpload {
				diskStatus, err = sysinfo.GetDiskStatus(realPath)
				if err != nil {
					http.Error(w, "Failed to get disk info", http.StatusInternalServerError)
					return
				}
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
				http.Error(w, "Failed to execute template", http.StatusInternalServerError)
				return
			}
		} else { // file is not dir and file is served

			SetCustomHeaders(w, realPath)
			http.ServeFile(w, r, realPath)
		}
	}
}

// Set custom options in the header, because some files like .mkv are treated by the browser as .webm .
// Other special cases will be added in future.
func SetCustomHeaders(w http.ResponseWriter, path string) {

	ext := strings.ToLower(filepath.Ext(path))
	filename := filepath.Base(path)

	switch ext {
	case ".mkv":
		w.Header().Set("Content-Type", "video/x-matroska")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	// Add other after testing
	default:
		// Do nothing; let http.ServeFile handle it
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
