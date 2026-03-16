package services

import (
	"context"
	"testing"
	"time"
	"uptime-checker/internal/api/dto"
	"uptime-checker/internal/api/models"
	sharedDto "uptime-checker/internal/shared/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestSiteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Site{}, &models.SiteCheck{})
	require.NoError(t, err)

	return db
}

func TestSiteService_CreateSite(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	userID := uuid.New()
	req := dto.CreateSiteRequest{
		URL:         "https://example.com",
		Name:        "Example",
		IntervalSec: 60,
	}

	site, err := svc.CreateSite(req, userID)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, site.ID)
	assert.Equal(t, req.URL, site.URL)
	assert.Equal(t, req.Name, site.Name)
	assert.Equal(t, userID, site.UserID)
	assert.True(t, site.IsActive)
}

func TestSiteService_GetUserSites(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	userID := uuid.New()

	for i := 0; i < 15; i++ {
		db.Create(&models.Site{
			URL:      "https://example.com",
			Name:     "Site",
			UserID:   userID,
			IsActive: true,
		})
	}

	result, err := svc.GetUserSites(userID, 1, 10)

	require.NoError(t, err)
	assert.Len(t, result.Data, 10)
	assert.Equal(t, int64(15), result.TotalCount)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Size)
}

func TestSiteService_UpdateSite(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	userID := uuid.New()
	site := models.Site{
		URL:      "https://old.com",
		Name:     "Old",
		UserID:   userID,
		IsActive: true,
	}
	db.Create(&site)

	req := dto.UpdateSiteRequest{
		URL: ptr("https://new.com"),
	}

	updated, err := svc.UpdateSite(site.ID, userID, req)

	require.NoError(t, err)
	assert.Equal(t, "https://new.com", updated.URL)
	assert.Equal(t, "Old", updated.Name)
}

func TestSiteService_UpdateSite_NotFound(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	req := dto.UpdateSiteRequest{
		URL: ptr("https://new.com"),
	}

	_, err := svc.UpdateSite(uuid.New(), uuid.New(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSiteService_DeleteSite(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	userID := uuid.New()
	site := models.Site{
		URL:    "https://example.com",
		Name:   "ToDelete",
		UserID: userID,
	}
	db.Create(&site)

	err := svc.DeleteSite(site.ID, userID)

	require.NoError(t, err)

	var count int64
	db.Model(&models.Site{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestSiteService_DeleteSite_NotFound(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	err := svc.DeleteSite(uuid.New(), uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSiteService_HandleCheckResult(t *testing.T) {
	db := setupTestSiteDB(t)
	svc := &SiteService{DB: db}

	siteID := uuid.New()
	db.Create(&models.Site{
		URL:    "https://example.com",
		Name:   "Test",
		UserID: uuid.New(),
	})

	result := sharedDto.SiteCheckResult{
		SiteID:     siteID,
		StatusCode: 200,
		LatencyMs:  150,
		IsUp:       true,
		CheckedAt:  time.Now(),
	}

	err := svc.HandleCheckResult(context.Background(), result)

	require.NoError(t, err)

	var check models.SiteCheck
	err = db.First(&check, "site_id = ?", siteID).Error
	require.NoError(t, err)
	assert.Equal(t, 200, check.StatusCode)
	assert.True(t, check.IsUp)
}

func ptr[T any](v T) *T {
	return &v
}
