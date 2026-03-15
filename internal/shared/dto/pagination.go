package dto

import "math"

type PagedResponse[T any] struct {
	Data       T     `json:"data"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalPages int   `json:"total_pages"`
}

func NewPagedResponse[T any](data T, total int64, page, size int) *PagedResponse[T] {
	return &PagedResponse[T]{
		Data:       data,
		TotalCount: total,
		Page:       page,
		Size:       size,
		TotalPages: int(math.Ceil(float64(total) / float64(size))),
	}
}
