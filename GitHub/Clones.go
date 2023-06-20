package GitHub

import (
	"GitHub-Analytics-Tracker/Database"
	"GitHub-Analytics-Tracker/GitHub/GithubAuth"
	"context"
	"github.com/google/go-github/v53/github"
	"gorm.io/gorm"
	"sync"
)

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

func WillInsertCloneData(data []*Database.CloneInstance, db *gorm.DB) {
	result := db.Create(&data)
	for result.Error != nil {
		result = db.Create(&data)
	}
}
