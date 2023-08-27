package Database

import (
	DBUtil "GitHub-Analytics-Tracker/Database/Util"
	"GitHub-Analytics-Tracker/Util"
	"log"

	"gorm.io/gorm"
)

var (
	// dEV_TOKEN_ENV_KEY is the default name for the ENV key for the GitHub Access token
	dEV_TOKEN_ENV_KEY = "DEV_REPO_API_KEY"
)

// Repository represents a repository that must be tracked
type Repository struct {
	// Owner is the GitHub owner
	Owner string `gorm:"primaryKey"`
	// Repo is the name of the repository
	Repo string `gorm:"primaryKey"`
	// TokenKey represents the name(key) of the env variable where the api key for this repo is found
	TokenKey string
}

func GetAllRepositories(db *gorm.DB) (repositories []*Repository, err error) {
	result := db.Find(&repositories)
	return repositories, result.Error
}

// migrateRepository migrates the Repository type to the DB ensuring they keep parity
func migrateRepository(db *gorm.DB) error {
	return db.AutoMigrate(Repository{})
}

// mustCreateDevRepo tries to create the local development repo object and panics if it is unable to construct it
func mustCreateDevRepo() *Repository {
	owner, err := Util.GetENV("DEV_REPO_OWNER")
	if err != nil {
		panic(err)
	}
	name, err := Util.GetENV("DEV_REPO")
	if err != nil {
		panic(err)
	}

	return &Repository{
		Owner:    owner,
		Repo:     name,
		TokenKey: dEV_TOKEN_ENV_KEY,
	}
}

// mustSetupDevEnv tries to set up the dev database for development, panics if it is unable to do so
func mustSetupDevEnv(db *gorm.DB) {
	err := migrateRepository(db)
	if err != nil {
		panic(err)
	}
	mustCleanUpDevTables(db)
	mustLoadDevRepo(db)
}

// mustCleanUpDevRepo deletes the devRepo record from the db and panics if it is unable to do so
func mustCleanUpDevRepo(devRepo *Repository, db *gorm.DB) {
	Result := db.Delete(devRepo)
	if Result.Error != nil {
		log.Panic(Result.Error)
	}
}

func mustCleanUpDevTables(db *gorm.DB) {
	session := db.Session(&gorm.Session{AllowGlobalUpdate: true})
	Result := session.Delete(&CloneInstance{})
	if Result.Error != nil {
		log.Panic(Result.Error)
	}
	Result = session.Delete(&ReferralPath{})
	if Result.Error != nil {
		log.Panic(Result.Error)
	}
	Result = session.Delete(&ReferralSource{})
	if Result.Error != nil {
		log.Panic(Result.Error)
	}
}

// mustLoadDevRepo tries to load the dev repo into the database and panics if it is unable to
func mustLoadDevRepo(db *gorm.DB) {
	DevRepo := mustCreateDevRepo()
	mustCleanUpDevRepo(DevRepo, db)
	Result := db.Create(DevRepo)
	if Result.Error != nil {
		log.Panic(Result.Error)
	}
}

// MustSetupDB sets up the database so that it is ready for use and panics if it cannot
func MustSetupDB(stage Util.Stage) *gorm.DB {
	db := DBUtil.MustConnect()

	//TODO make this neater and less copy pasty
	err := migrateRepository(db)
	if err != nil {
		panic(err)
	}
	err = migrateCloneInstance(db)
	if err != nil {
		panic(err)
	}

	err = migrateReferralSource(db)
	if err != nil {
		panic(err)
	}

	err = migrateReferralPaths(db)
	if err != nil {
		panic(err)
	}

	switch stage {
	case Util.Dev:
		mustSetupDevEnv(db)
	}
	return db
}
