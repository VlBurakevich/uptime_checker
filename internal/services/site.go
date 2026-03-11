package services

import (
	"errors"
	"uptime-checker/internal/dto"
	"uptime-checker/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SiteService struct {
	DB *gorm.DB
}

func (s *SiteService) GetUserSites(userId uuid.UUID) ([]models.Site, error) {
	var sites []models.Site
	err := s.DB.Where("user_id = ?", userId).Find(&sites).Error
	return sites, err
}

func (s *SiteService) CreateSite(req dto.CreateSiteRequest, userId uuid.UUID) (*models.Site, error) {
	site := models.Site{
		URL:      req.URL,
		Name:     req.Name,
		Interval: req.Interval,
		UserID:   userId,
		IsActive: true,
	}

	if err := s.DB.Create(&site).Error; err != nil {
		return nil, err
	}
	return &site, nil
}

func (s *SiteService) UpdateSite(siteID uuid.UUID, userID uuid.UUID, req dto.UpdateSiteRequest) (*models.Site, error) {
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
	if req.Interval != nil {
		site.Interval = *req.Interval
	}

	if err := s.DB.Save(&site).Error; err != nil {
		return nil, err
	}

	return &site, nil
}
