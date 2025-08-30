package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	storageDir        = "storage"
	indexesDir        = "_indexes"
	byBundleIDDir     = "by_bundle_id"
	buildInfoFileName = "build_info.json"
)

// IndexEntry represents an entry in the bundle ID index.
type IndexEntry struct {
	UploadID  string    `json:"upload_id"`
	CreatedAt time.Time `json:"created_at"`
}

// initStorage initializes the storage directories.
func initStorage() error {
	for _, dir := range []string{storageDir, filepath.Join(storageDir, indexesDir), filepath.Join(storageDir, indexesDir, byBundleIDDir)} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// getBuildInfo loads the build metadata from a file.
func getBuildInfo(uploadID string) (*BuildInfo, error) {
	filePath := filepath.Join(storageDir, uploadID, buildInfoFileName)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("build info not found for upload ID %s", uploadID)
		}
		return nil, fmt.Errorf("failed to open build info file: %w", err)
	}
	defer file.Close()

	var info BuildInfo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode build info: %w", err)
	}

	return &info, nil
}

// getIndexEntriesForBundleID returns the index entries for a given bundle ID.
func getIndexEntriesForBundleID(bundleID string) ([]IndexEntry, error) {
	indexFilePath := filepath.Join(storageDir, indexesDir, byBundleIDDir, fmt.Sprintf("%s.json", bundleID))

	file, err := os.Open(indexFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no versions found for bundle ID %s", bundleID)
		}
		return nil, fmt.Errorf("failed to open index file: %w", err)
	}
	defer file.Close()

	var index []IndexEntry
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&index); err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to decode index file: %w", err)
	}

	return index, nil
}

// getLatestUploadIDForBundleID returns the upload ID of the latest version for a given bundle ID.
func getLatestUploadIDForBundleID(bundleID string) (string, error) {
	index, err := getIndexEntriesForBundleID(bundleID)
	if err != nil {
		return "", err
	}
	if len(index) == 0 {
		return "", fmt.Errorf("no versions found for bundle ID %s", bundleID)
	}
	return index[0].UploadID, nil
}

// getAllVersions returns all versions for a given bundle ID.
func getAllVersions(bundleID string) ([]*BuildInfo, error) {
	index, err := getIndexEntriesForBundleID(bundleID)
	if err != nil {
		return nil, err
	}

	var builds []*BuildInfo
	for _, entry := range index {
		build, err := getBuildInfo(entry.UploadID)
		if err != nil {
			// Log the error but continue, so one corrupted build doesn't fail the whole request
			fmt.Printf("Error getting build info for %s: %v\n", entry.UploadID, err)
			continue
		}
		builds = append(builds, build)
	}

	return builds, nil
}

// getLatestVersion returns the latest version for a given bundle ID.
func getLatestVersion(bundleID string) (*BuildInfo, error) {
	uploadID, err := getLatestUploadIDForBundleID(bundleID)
	if err != nil {
		return nil, err
	}
	return getBuildInfo(uploadID)
}

// getAllApps returns the latest version of all apps.
func getAllApps() ([]*BuildInfo, error) {
	allApps := make([]*BuildInfo, 0)
	bundleIDFiles, err := os.ReadDir(filepath.Join(storageDir, indexesDir, byBundleIDDir))
	if err != nil {
		if os.IsNotExist(err) {
			return allApps, nil // Return empty slice if directory doesn't exist
		}
		return nil, fmt.Errorf("failed to read bundle ID index directory: %w", err)
	}

	for _, file := range bundleIDFiles {
		if file.IsDir() {
			continue
		}
		bundleID := file.Name()
		bundleID = bundleID[:len(bundleID)-len(".json")]
		build, err := getLatestVersion(bundleID)
		if err != nil {
			fmt.Printf("Error getting latest version for %s: %v\n", bundleID, err)
			continue
		}
		allApps = append(allApps, build)
	}

	return allApps, nil
}

// saveBuildInfo saves the build metadata to a file.
func saveBuildInfo(info *BuildInfo) error {
	uploadDir := filepath.Join(storageDir, info.UploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(uploadDir, buildInfoFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create build info file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		return fmt.Errorf("failed to encode build info: %w", err)
	}

	return nil
}

// saveAppFile saves the application file.
func saveAppFile(info *BuildInfo, appFile io.Reader) error {
	uploadDir := filepath.Join(storageDir, info.UploadID)
	fileName := "app.ipa"
	if info.Platform == Android {
		fileName = "app.apk"
	}
	filePath := filepath.Join(uploadDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create app file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, appFile)
	if err != nil {
		return fmt.Errorf("failed to save app file: %w", err)
	}

	return nil
}

// updateIndex adds a new entry to the bundle ID index.
func updateIndex(info *BuildInfo) error {
	indexFilePath := filepath.Join(storageDir, indexesDir, byBundleIDDir, fmt.Sprintf("%s.json", info.BundleID))

	var index []IndexEntry
	file, err := os.Open(indexFilePath)
	if err == nil {
		// Index file exists, read it
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&index); err != nil && err != io.EOF {
			file.Close()
			return fmt.Errorf("failed to decode index file: %w", err)
		}
		file.Close()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to open index file: %w", err)
	}

	// Add new entry
	index = append(index, IndexEntry{
		UploadID:  info.UploadID,
		CreatedAt: info.CreatedAt,
	})

	// Sort by creation date to keep it organized
	sort.Slice(index, func(i, j int) bool {
		return index[i].CreatedAt.After(index[j].CreatedAt)
	})

	// Write back to the index file
	file, err = os.Create(indexFilePath)
	if err != nil {
		return fmt.Errorf("failed to create index file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(index); err != nil {
		return fmt.Errorf("failed to encode index: %w", err)
	}

	return nil
}

// SaveUpload saves the build info, app file, and updates the index.
func SaveUpload(info *BuildInfo, appFile io.Reader) error {
	if err := saveBuildInfo(info); err != nil {
		return err
	}
	if err := saveAppFile(info, appFile); err != nil {
		return err
	}
	if err := updateIndex(info); err != nil {
		return err
	}
	return nil
}
