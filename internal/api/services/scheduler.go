package services

import (
	"context"
	"time"
	"uptime-checker/internal/api/broker"
	"uptime-checker/internal/api/models"
	"uptime-checker/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Scheduler struct {
	db       *gorm.DB
	producer *broker.TaskProducer
	limit    int
}

func NewScheduler(db *gorm.DB, producer *broker.TaskProducer, limit int) *Scheduler {
	return &Scheduler{
		db:       db,
		producer: producer,
		limit:    limit,
	}
}

func (s *Scheduler) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.process(ctx)
		}
	}
}

func (s *Scheduler) process(ctx context.Context) {
	var sites []models.Site

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
			Where("is_active = ? AND (last_checked_at IS NULL OR last_checked_at + (interval * interval '1 second') <= NOW())", true).
			Limit(s.limit).
			Find(&sites).Error

		if err != nil || len(sites) == 0 {
			return err
		}

		siteIDs := make([]uuid.UUID, len(sites))
		for i, site := range sites {
			siteIDs[i] = site.ID
		}

		return tx.Model(&models.Site{}).Where("id IN ?", siteIDs).Update("last_checked_at", time.Now()).Error
	})

	if err != nil || len(sites) == 0 {
		return
	}

	for _, site := range sites {
		task := dto.SiteCheckTask{
			SiteID: site.ID,
			URL:    site.URL,
		}
		_ = s.producer.PublishTask(ctx, task)
	}
}
