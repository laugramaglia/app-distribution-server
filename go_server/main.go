package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/base64"

	"github.com/google/uuid"
	"github.com/nao1215/deapk/apk"
	"github.com/skip2/go-qrcode"
)

// DownloadResponse represents the response for the download endpoint.
type DownloadResponse struct {
	BuildInfo
	QRCode string `json:"qr_code"`
}

func main() {
	if err := initStorage(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	http.HandleFunc("/api/apps", appsHandler)
	http.HandleFunc("/api/apps/upload", uploadHandler)
	http.HandleFunc("/api/apps/", appVersionsHandler) // Note the trailing slash

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func appsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apps, err := getAllApps()
	if err != nil {
		http.Error(w, "Failed to get apps", http.StatusInternalServerError)
		log.Printf("Error getting apps: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding apps: %v", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB limit
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("app_file")
	if err != nil {
		http.Error(w, "Failed to get app file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var platform Platform
	var buildInfo BuildInfo

	if strings.HasSuffix(handler.Filename, ".apk") {
		platform = Android

		// Create a temporary file to store the uploaded apk
		tmpfile, err := os.CreateTemp("", "upload-*.apk")
		if err != nil {
			http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name()) // clean up

		// Copy the uploaded file to the temporary file
		fileSize, err := io.Copy(tmpfile, file)
		if err != nil {
			http.Error(w, "Failed to save temporary file", http.StatusInternalServerError)
			return
		}

		// Parse the APK file
		apkParser := apk.NewAPK(tmpfile.Name())
		if err := apkParser.Parse(); err != nil {
			http.Error(w, "Failed to parse apk file", http.StatusInternalServerError)
			log.Printf("Error parsing apk: %v", err)
			return
		}

		buildInfo = BuildInfo{
			UploadID:    uuid.New().String(),
			BundleID:    apkParser.Package.Basic.PackageName,
			Version:     apkParser.Package.Basic.Version,
			BuildNumber: "0", // This library does not seem to extract the build number.
			Title:       apkParser.Package.Basic.ApplicationName,
			FileSize:    fileSize,
			CreatedAt:   time.Now(),
			Platform:    platform,
		}

		// Reset the file reader to the beginning
		if _, err := tmpfile.Seek(0, 0); err != nil {
			http.Error(w, "Failed to seek temporary file", http.StatusInternalServerError)
			return
		}

		if err := SaveUpload(&buildInfo, tmpfile); err != nil {
			http.Error(w, "Failed to save upload", http.StatusInternalServerError)
			log.Printf("Error saving upload: %v", err)
			return
		}

	} else if strings.HasSuffix(handler.Filename, ".ipa") {
		platform = IOS

		bundleID := r.FormValue("bundle_id")
		version := r.FormValue("version")
		buildNumber := r.FormValue("build_number")
		title := r.FormValue("title")

		if bundleID == "" || version == "" || buildNumber == "" || title == "" {
			http.Error(w, "Missing required metadata for .ipa upload (bundle_id, version, build_number, title)", http.StatusBadRequest)
			return
		}

		// Create a temporary file to get the file size
		tmpfile, err := os.CreateTemp("", "upload-*.ipa")
		if err != nil {
			http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name()) // clean up

		fileSize, err := io.Copy(tmpfile, file)
		if err != nil {
			http.Error(w, "Failed to save temporary file", http.StatusInternalServerError)
			return
		}

		buildInfo = BuildInfo{
			UploadID:    uuid.New().String(),
			BundleID:    bundleID,
			Version:     version,
			BuildNumber: buildNumber,
			Title:       title,
			FileSize:    fileSize,
			CreatedAt:   time.Now(),
			Platform:    platform,
		}

		// Reset the file reader to the beginning
		if _, err := tmpfile.Seek(0, 0); err != nil {
			http.Error(w, "Failed to seek temporary file", http.StatusInternalServerError)
			return
		}


		if err := SaveUpload(&buildInfo, tmpfile); err != nil {
			http.Error(w, "Failed to save upload", http.StatusInternalServerError)
			log.Printf("Error saving upload: %v", err)
			return
		}

	} else {
		http.Error(w, "Invalid file type. Only .apk and .ipa files are supported", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(buildInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func appVersionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	bundleID := parts[3]

	// This is the /api/apps/:id/versions endpoint
	if len(parts) > 4 && parts[4] == "versions" {
		versions, err := getAllVersions(bundleID)
		if err != nil {
			http.Error(w, "Failed to get versions", http.StatusInternalServerError)
			log.Printf("Error getting versions for %s: %v", bundleID, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(versions); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding versions: %v", err)
		}
		return
	}

	// This is the /api/apps/:id/download endpoint, which is the default
	// for /api/apps/:id
	build, err := getLatestVersion(bundleID)
	if err != nil {
		http.Error(w, "Failed to get latest version", http.StatusInternalServerError)
		log.Printf("Error getting latest version for %s: %v", bundleID, err)
		return
	}

	downloadURL := "http://" + r.Host + "/api/apps/" + bundleID + "/download" // This should point to the actual file download, which is not implemented yet. For now, it's a placeholder.
	var png []byte
	png, err = qrcode.Encode(downloadURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		log.Printf("Error generating QR code: %v", err)
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(png)

	response := DownloadResponse{
		BuildInfo: *build,
		QRCode:    qrCodeBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding download response: %v", err)
	}
}
