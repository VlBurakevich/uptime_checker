package database

import "gorm.io/gorm"

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}
		if size > 100 {
			size = 100
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
