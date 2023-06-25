package main

import (
	"GitHub-Analytics-Tracker/Database"
	"GitHub-Analytics-Tracker/GitHub"
	"GitHub-Analytics-Tracker/GitHub/GithubAuth"
	"GitHub-Analytics-Tracker/Util"
	"fmt"
	"sync"
)

var (
	master   sync.WaitGroup
	finished = make(chan bool, 5)
)

func main() {
	stage := Util.MustSetupEnvironment()
	db := Database.MustSetupDB(stage)
	go GithubAuth.ValidateRepositories(db)
	repos, err := Database.GetAllRepositories(db)
	for err != nil {
		//TODO better error handling
		repos, err = Database.GetAllRepositories(db)
	}
	master.Add(len(repos) * 3)
	go waitForComplete(&master)

L:
	for {
		select {
		case <-finished:
			break L
		case connection := <-GithubAuth.Connections:
			//start a separate routine for each data type
			go GitHub.WillGetClones(connection, &master, db)
			go GitHub.WillGetTrafficReferrers(connection, &master, db)
			go GitHub.WillGetTrafficPaths(connection, &master, db)
		}
	}
	fmt.Println("GitHub Analytics Tracker is Finished")

}

// waitForComplete is intended to run as a rp
func waitForComplete(group *sync.WaitGroup) {
	group.Wait()
	finished <- true
}
