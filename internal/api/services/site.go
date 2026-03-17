package services

import (
	"context"
	"errors"
	"uptime-checker/internal/api/database"
	apiDto "uptime-checker/internal/api/dto"
	"uptime-checker/internal/api/models"
	sharedDto "uptime-checker/internal/shared/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SiteService struct {
	DB *gorm.DB
}

func (s *SiteService) GetUserSite(userID uuid.UUID, siteID uuid.UUID) (*models.Site, error) {
	site := &models.Site{}

	if err := s.DB.Model(&models.Site{}).Where("id = ? AND user_id = ?", siteID, userID).First(site).Error; err != nil {
		return nil, errors.New("site not found")
	}
	return site, nil
}

func (s *SiteService) GetUserSites(userId uuid.UUID, page, size int) (*sharedDto.PagedResponse[[]models.Site], error) {
	var sites []models.Site
	var total int64

	if err := s.DB.Model(&models.Site{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, err
	}

	err := s.DB.Where("user_id = ?", userId).
		Scopes(database.Paginate(page, size)).
		Find(&sites).Error

	if err != nil {
		return nil, err
	}

	return sharedDto.NewPagedResponse(sites, total, page, size), nil
}

func (s *SiteService) GetUserSiteChecks(userID uuid.UUID, page int, size int) (*sharedDto.PagedResponse[[]models.SiteCheck], error) {
	var site models.Site
	if err := s.DB.Where("user_id = ?", userID).First(&site).Error; err != nil {
		return nil, errors.New("site not found")
	}

	var siteChecks []models.SiteCheck
	var total int64

	if err := s.DB.Model(&models.SiteCheck{}).Where("site_id = ?", site.ID).Count(&total).Error; err != nil {
		return nil, err
	}

	err := s.DB.Where("site_id", site.ID).
		Scopes(database.Paginate(page, size)).
		Find(&siteChecks).Error

	if err != nil {
		return nil, err
	}

	return sharedDto.NewPagedResponse(siteChecks, total, page, size), nil
}

func (s *SiteService) CreateSite(req apiDto.CreateSiteRequest, userId uuid.UUID) (*models.Site, error) {
	site := models.Site{
		URL:         req.URL,
		Name:        req.Name,
		IntervalSec: req.IntervalSec,
		UserID:      userId,
		IsActive:    true,
	}

	if err := s.DB.Create(&site).Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (s *SiteService) UpdateSite(siteID uuid.UUID, userID uuid.UUID, req apiDto.UpdateSiteRequest) (*models.Site, error) {
	var site models.Site
	if err := s.DB.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		return nil, errors.New("site not found or access denied")
	}

	if req.URL != nil {
		site.URL = *req.URL
	}
	if req.Name != nil {
		site.Name = *req.Name
	}
	if req.IntervalSec != nil {
		site.IntervalSec = *req.IntervalSec
	}

	if err := s.DB.Save(&site).Error; err != nil {
		return nil, err
	}

	return &site, nil
}

func (s *SiteService) DeleteSite(siteID uuid.UUID, userID uuid.UUID) error {
	result := s.DB.Where("id = ? AND user_id = ?", siteID, userID).Delete(&models.Site{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("site not found or access denied")
	}

	return nil
}

func (s *SiteService) HandleCheckResult(ctx context.Context, res sharedDto.SiteCheckResult) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		check := models.SiteCheck{
			SiteID:     res.SiteID,
			StatusCode: res.StatusCode,
			LatencyMs:  res.LatencyMs,
			IsUp:       res.IsUp,
			CheckedAt:  res.CheckedAt,
		}

		if err := tx.Create(&check).Error; err != nil {
			return err
		}

		return nil
	})
}
