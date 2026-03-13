package dto

type CreateSiteRequest struct {
	URL         string `json:"url" binding:"required,url"`
	Name        string `json:"name" binding:"required,min=2"`
	IntervalSec int    `json:"interval_sec" binding:"required,min=1"`
}

type UpdateSiteRequest struct {
	URL         *string `json:"url" binding:"omitempty,url"`
	Name        *string `json:"name" binding:"omitempty,min=2"`
	IntervalSec *int    `json:"interval_sec" binding:"omitempty,min=1"`
}
