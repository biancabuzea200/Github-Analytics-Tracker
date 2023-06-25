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

// WillGetTrafficReferrers will retrieve the referral paths for the supplied github repository and will retry until successful (TODO better retry logic)
func WillGetTrafficReferrers(connection *GithubAuth.Connection, wg *sync.WaitGroup, db *gorm.DB) {
	defer wg.Done()
	var referralSources []*Database.ReferralSource

	ctx := context.Background()
	referrals, _, err := connection.Client.Repositories.ListTrafficReferrers(ctx, connection.Repository.Owner, connection.Repository.Repo)
	//TODO some sort of better retry logic here
	for err != nil {
		referrals, _, err = connection.Client.Repositories.ListTrafficReferrers(ctx, connection.Repository.Owner, connection.Repository.Repo)
	}

	today := time.Now()

	for _, referrer := range referrals {
		referralSources = append(referralSources, &Database.ReferralSource{
			Repository: connection.Repository,
			Day:        &today,
			Referrer:   referrer.Referrer,
			Count:      referrer.Count,
			Uniques:    referrer.Uniques,
		},
		)
	}
	WillInsertSources(referralSources, db)
}

// WillInsertSources is guaranteed to make the referralPaths database insert  will get stuck in an infinite loop
func WillInsertSources(data []*Database.ReferralSource, db *gorm.DB) {
	result := db.Create(&data)
	//TODO better error handling
	for result.Error != nil {
		result = db.Create(&data)
	}
}
