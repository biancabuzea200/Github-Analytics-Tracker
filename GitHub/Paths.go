package GitHub

import (
	"GitHub-Analytics-Tracker/Database"
	"GitHub-Analytics-Tracker/GitHub/GithubAuth"
	"context"
	"gorm.io/gorm"
	"sync"
	"time"
)

//TODO very similar logic, could eventually be pulled out but this would cause unnecessary coupling

// WillGetTrafficPaths will retrieve the referral sources for the supplied github repository and will retry until successful (TODO better retry logic)
func WillGetTrafficPaths(connection *GithubAuth.Connection, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	var trafficPaths []*Database.ReferralPath

	ctx := context.Background()
	paths, _, err := connection.Client.Repositories.ListTrafficPaths(ctx, connection.Repository.Owner, connection.Repository.Repo)
	//TODO some sort of better retry logic here
	for err != nil {
		paths, _, err = connection.Client.Repositories.ListTrafficPaths(ctx, connection.Repository.Owner, connection.Repository.Repo)
	}

	today := time.Now()

	for _, path := range paths {
		trafficPaths = append(trafficPaths, &Database.ReferralPath{
			Repository:  connection.Repository,
			Day:         &today,
			TrafficPath: path,
		},
		)
	}
	WillInsertPaths(trafficPaths, db)
}

// WillInsertPaths is guaranteed to make the referralPaths database insert  will get stuck in an infinite loop
func WillInsertPaths(data []*Database.ReferralPath, db *gorm.DB) {
	result := db.Create(&data)
	//TODO better error handling
	for result.Error != nil {
		result = db.Create(&data)
	}
}
