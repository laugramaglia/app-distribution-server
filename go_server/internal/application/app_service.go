package application

import (
	"app-distribution-server-go/internal/domain"
	"io"
)

type AppRepository interface {
	GetAllApps() ([]*domain.BuildInfo, error)
	GetAllVersions(bundleID string) ([]*domain.BuildInfo, error)
	GetLatestVersion(bundleID string) (*domain.BuildInfo, error)
	GetBuild(bundleID, version, buildNumber string) (*domain.BuildInfo, error)
	SaveUpload(info *domain.BuildInfo, appFile io.Reader) error
}

type AppService struct {
	repo AppRepository
}

func NewAppService(repo AppRepository) *AppService {
	return &AppService{repo: repo}
}

func (s *AppService) GetAllApps() ([]*domain.BuildInfo, error) {
	return s.repo.GetAllApps()
}

func (s *AppService) GetLatestVersion(bundleID string) (*domain.BuildInfo, error) {
	return s.repo.GetLatestVersion(bundleID)
}

func (s *AppService) GetAllVersions(bundleID string) ([]*domain.BuildInfo, error) {
	return s.repo.GetAllVersions(bundleID)
}

func (s *AppService) GetBuild(bundleID, version, buildNumber string) (*domain.BuildInfo, error) {
	return s.repo.GetBuild(bundleID, version, buildNumber)
}

func (s *AppService) SaveUpload(info *domain.BuildInfo, appFile io.Reader) error {
	return s.repo.SaveUpload(info, appFile)
}
