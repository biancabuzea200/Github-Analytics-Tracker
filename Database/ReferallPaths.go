package Database

import (
	"github.com/google/go-github/v53/github"
	"gorm.io/gorm"
	"time"
)

type ReferralSource struct {
	RepositoryOwner string
	RepositoryRepo  string
	Repository      *Repository `gorm:"primaryKey"`
	Day             *time.Time  `gorm:"primaryKey"`
	Referrer        *string
	Count           *int
	Uniques         *int
}

// migrateReferralSource migrates the ReferralSource struct as a table in the DB
func migrateReferralSource(db *gorm.DB) error {
	return db.AutoMigrate(ReferralSource{})
}

type ReferralPath struct {
	RepositoryOwner string
	RepositoryRepo  string
	Repository      *Repository         `gorm:"primaryKey"`
	Day             *time.Time          `gorm:"primaryKey"`
	TrafficPath     *github.TrafficPath `gorm:"embedded"`
}

// migrateReferralPaths migrates the ReferralPath struct as a table in the DB
func migrateReferralPaths(db *gorm.DB) error {
	return db.AutoMigrate(ReferralPath{})
}
