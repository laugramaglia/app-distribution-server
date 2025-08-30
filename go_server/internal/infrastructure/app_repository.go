package infrastructure

import (
	"app-distribution-server-go/internal/domain"
	"io"
)

type FileAppRepository struct{}

func NewFileAppRepository() *FileAppRepository {
	return &FileAppRepository{}
}

func (r *FileAppRepository) GetAllApps() ([]*domain.BuildInfo, error) {
	return GetAllApps()
}

func (r *FileAppRepository) GetAllVersions(bundleID string) ([]*domain.BuildInfo, error) {
	return GetAllVersions(bundleID)
}

func (r *FileAppRepository) GetLatestVersion(bundleID string) (*domain.BuildInfo, error) {
	return GetLatestVersion(bundleID)
}

func (r *FileAppRepository) SaveUpload(info *domain.BuildInfo, appFile io.Reader) error {
	return SaveUpload(info, appFile)
}
