package Database

import (
	"gorm.io/gorm"
	"time"
)

type CloneInstance struct {
	RepositoryOwner string
	RepositoryRepo  string
	Repository      *Repository `gorm:"primaryKey"`
	Day             *time.Time  `gorm:"primaryKey"`
	Clones          *int
	Uniques         *int
}

func migrateCloneInstance(db *gorm.DB) error {
	return db.AutoMigrate(CloneInstance{})
}
