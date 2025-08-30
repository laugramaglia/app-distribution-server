package interfaces

import (
	"app-distribution-server-go/internal/application"
	"app-distribution-server-go/internal/domain"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nao1215/deapk/apk"
	"github.com/skip2/go-qrcode"
)

type AppHandlers struct {
	service *application.AppService
}

func NewAppHandlers(service *application.AppService) *AppHandlers {
	return &AppHandlers{service: service}
}

// DownloadResponse represents the response for the download endpoint.
type DownloadResponse struct {
	domain.BuildInfo
	QRCode string `json:"qr_code"`
}

// AppsHandler godoc
// @Summary List all apps
// @Description Get a list of all available applications.
// @Tags apps
// @Produce  json
// @Success 200 {array} domain.BuildInfo
// @Failure 500 {string} string "Failed to get apps"
// @Router /apps [get]
func (h *AppHandlers) AppsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apps, err := h.service.GetAllApps()
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

// UploadHandler godoc
// @Summary Upload a new app
// @Description Upload a new .apk or .ipa file.
// @Tags apps
// @Accept  multipart/form-data
// @Produce  json
// @Param   app_file formData file true  "Application file (.apk or .ipa)"
// @Param   bundle_id formData string false "Bundle ID (for .ipa)"
// @Param   version formData string false "Version (for .ipa)"
// @Param   build_number formData string false "Build Number (for .ipa)"
// @Param   title formData string false "Title (for .ipa)"
// @Success 200 {object} domain.BuildInfo
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /apps/upload [post]
func (h *AppHandlers) UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	var platform domain.Platform
	var buildInfo domain.BuildInfo

	if strings.HasSuffix(handler.Filename, ".apk") {
		platform = domain.Android

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

		buildInfo = domain.BuildInfo{
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

		if err := h.service.SaveUpload(&buildInfo, tmpfile); err != nil {
			http.Error(w, "Failed to save upload", http.StatusInternalServerError)
			log.Printf("Error saving upload: %v", err)
			return
		}

	} else if strings.HasSuffix(handler.Filename, ".ipa") {
		platform = domain.IOS

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

		buildInfo = domain.BuildInfo{
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


		if err := h.service.SaveUpload(&buildInfo, tmpfile); err != nil {
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

// GetLatestAppVersionHandler godoc
// @Summary Get latest app version
// @Description Get a download link for the latest version of an app.
// @Tags apps
// @Produce  json
// @Param   bundle_id path string true "Bundle ID of the app"
// @Success 200 {object} DownloadResponse
// @Failure 400 {string} string "Invalid URL"
// @Failure 500 {string} string "Internal Server Error"
// @Router /apps/{bundle_id} [get]
func (h *AppHandlers) GetLatestAppVersionHandler(w http.ResponseWriter, r *http.Request) {
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

	build, err := h.service.GetLatestVersion(bundleID)
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

// GetAllAppVersionsHandler godoc
// @Summary Get all app versions
// @Description Get all versions of an app.
// @Tags apps
// @Produce  json
// @Param   bundle_id path string true "Bundle ID of the app"
// @Success 200 {array} domain.BuildInfo
// @Failure 400 {string} string "Invalid URL"
// @Failure 500 {string} string "Internal Server Error"
// @Router /apps/{bundle_id}/versions [get]
func (h *AppHandlers) GetAllAppVersionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	bundleID := parts[3]

	versions, err := h.service.GetAllVersions(bundleID)
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
}
