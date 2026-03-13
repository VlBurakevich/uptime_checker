package services

import (
	"context"
	"time"
	"uptime-checker/internal/api/broker"
	"uptime-checker/internal/api/models"
	"uptime-checker/internal/shared/dto"

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
		var siteIDs []uuid.UUID
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
			Model(&models.Site{}).
			Where("is_active = ? AND next_checked_at <= NOW()", true).
			Limit(s.limit).
			Pluck("id", &siteIDs).Error

		if err != nil || len(siteIDs) == 0 {
			return err
		}

		if err := tx.Where("id IN ?", siteIDs).Find(&sites).Error; err != nil {
			return err
		}

		return tx.Model(&models.Site{}).
			Where("id IN ?", siteIDs).
			Updates(map[string]interface{}{
				"last_checked_at": gorm.Expr("NOW()"),
				"next_checked_at": gorm.Expr("GREATEST(next_checked_at + (interval_sec * '1 second'::interval), NOW() + (interval_sec * '1 second'::interval))"),
			}).Error
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
