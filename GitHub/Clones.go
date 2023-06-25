package GitHub

import (
	"context"
	"sync"

	"github.com/google/go-github/v53/github"
	"gorm.io/gorm"

	"GitHub-Analytics-Tracker/Database"
	"GitHub-Analytics-Tracker/GitHub/GithubAuth"
)

// WillGetClones will retireve the clones for the supples github repository and will retry until successful (TODO better retry logic)
func WillGetClones(connection *GithubAuth.Connection, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	var cloneData []*Database.CloneInstance

	ctx := context.Background()
	opt := github.TrafficBreakdownOptions{Per: "day"}
	clones, _, err := connection.Client.Repositories.ListTrafficClones(ctx, connection.Repository.Owner, connection.Repository.Repo, &opt)
	//TODO some sort of better retry logic here
	for err != nil {
		clones, _, err = connection.Client.Repositories.ListTrafficClones(ctx, connection.Repository.Owner, connection.Repository.Repo, &opt)
	}
	for _, day := range clones.Clones {
		cloneData = append(cloneData, &Database.CloneInstance{
			Repository: connection.Repository,
			Day:        day.Timestamp.GetTime(),
			Clones:     day.Count,
			Uniques:    day.Uniques,
		})
	}
	WillInsertCloneData(cloneData, db)
}

// WillInsertCloneData is garuanteed to make the Database.CloneInstance database insert or it will get stuck in an infinite loop
func WillInsertCloneData(data []*Database.CloneInstance, db *gorm.DB) {
	result := db.Create(&data)
	//TODO better error handling
	for result.Error != nil {
		result = db.Create(&data)
	}
}
