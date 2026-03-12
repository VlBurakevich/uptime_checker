package handlers

import (
	"net/http"
	"uptime-checker/internal/api/dto"
	"uptime-checker/internal/api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SiteHandler struct {
	Service *services.SiteService
}

func (h *SiteHandler) List(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	sites, err := h.Service.GetUserSites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sites)
}

func (h *SiteHandler) Update(c *gin.Context) {
	var req dto.UpdateSiteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.URL == nil && req.Name == nil && req.Interval == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	siteID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid site id"})
		return
	}

	site, err := h.Service.UpdateSite(siteID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, site)
}

func (h *SiteHandler) Create(c *gin.Context) {
	var req dto.CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	site, err := h.Service.CreateSite(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, site)
}

func (h *SiteHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
	id, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	uid, ok := id.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return uuid.Nil, false
	}
	return uid, true
}
