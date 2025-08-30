package main

import "time"

// Platform represents the mobile platform (iOS or Android).
type Platform string

const (
	// IOS is the iOS platform.
	IOS Platform = "ios"
	// Android is the Android platform.
	Android Platform = "android"
)

// BuildInfo represents the metadata for a single build of an application.
type BuildInfo struct {
	UploadID    string    `json:"upload_id"`
	BundleID    string    `json:"bundle_id"`
	Version     string    `json:"version"`
	BuildNumber string    `json:"build_number"`
	Title       string    `json:"title"`
	Icon        string    `json:"icon,omitempty"`
	Description string    `json:"description,omitempty"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	Platform    Platform  `json:"platform"`
}
