package infrastructure

import (
	"app-distribution-server-go/internal/domain"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PostgresAppRepository struct {
	db *sql.DB
}

func NewPostgresAppRepository(db *sql.DB) (*PostgresAppRepository, error) {
	return &PostgresAppRepository{db: db}, nil
}

func (r *PostgresAppRepository) GetAllApps() ([]*domain.BuildInfo, error) {
	query := `
		SELECT upload_id, bundle_id, version, build_number, title, icon, description, file_size, created_at, platform
		FROM (
			SELECT *, ROW_NUMBER() OVER(PARTITION BY bundle_id ORDER BY created_at DESC) as rn
			FROM builds
		) t
		WHERE rn = 1
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query for all apps: %w", err)
	}
	defer rows.Close()

	var allApps []*domain.BuildInfo
	for rows.Next() {
		var app domain.BuildInfo
		if err := rows.Scan(&app.UploadID, &app.BundleID, &app.Version, &app.BuildNumber, &app.Title, &app.Icon, &app.Description, &app.FileSize, &app.CreatedAt, &app.Platform); err != nil {
			return nil, fmt.Errorf("failed to scan app row: %w", err)
		}
		allApps = append(allApps, &app)
	}

	return allApps, nil
}

func (r *PostgresAppRepository) GetAllVersions(bundleID string) ([]*domain.BuildInfo, error) {
	query := `
		SELECT upload_id, bundle_id, version, build_number, title, icon, description, file_size, created_at, platform
		FROM builds
		WHERE bundle_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, bundleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query for all versions of app %s: %w", bundleID, err)
	}
	defer rows.Close()

	var builds []*domain.BuildInfo
	for rows.Next() {
		var build domain.BuildInfo
		if err := rows.Scan(&build.UploadID, &build.BundleID, &build.Version, &build.BuildNumber, &build.Title, &build.Icon, &build.Description, &build.FileSize, &build.CreatedAt, &build.Platform); err != nil {
			return nil, fmt.Errorf("failed to scan build row: %w", err)
		}
		builds = append(builds, &build)
	}

	return builds, nil
}

func (r *PostgresAppRepository) GetLatestVersion(bundleID string) (*domain.BuildInfo, error) {
	query := `
		SELECT upload_id, bundle_id, version, build_number, title, icon, description, file_size, created_at, platform
		FROM builds
		WHERE bundle_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRow(query, bundleID)

	var build domain.BuildInfo
	if err := row.Scan(&build.UploadID, &build.BundleID, &build.Version, &build.BuildNumber, &build.Title, &build.Icon, &build.Description, &build.FileSize, &build.CreatedAt, &build.Platform); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no versions found for bundle ID %s", bundleID)
		}
		return nil, fmt.Errorf("failed to scan latest version row: %w", err)
	}

	return &build, nil
}

func (r *PostgresAppRepository) SaveUpload(info *domain.BuildInfo, appFile io.Reader) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `
		INSERT INTO builds (upload_id, bundle_id, version, build_number, title, icon, description, file_size, created_at, platform)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(query, info.UploadID, info.BundleID, info.Version, info.BuildNumber, info.Title, info.Icon, info.Description, info.FileSize, info.CreatedAt, info.Platform)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert build info: %w", err)
	}

	if err := r.saveAppFile(info, appFile); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// saveAppFile saves the application file.
func (r *PostgresAppRepository) saveAppFile(info *domain.BuildInfo, appFile io.Reader) error {
	uploadDir := filepath.Join(StorageDir, info.BundleID, info.UploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %w", err)
	}

	fileName := "app.ipa"
	if info.Platform == domain.Android {
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
