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
	finished chan bool
)

func main() {
	stage := Util.MustSetupEnvironment()
	db := Database.MustSetupDB(stage)
	GithubAuth.ValidateRepositories(db)

L:
	for {
		select {
		case <-finished:
			break L
		case connection := <-GithubAuth.Connections:
			fmt.Println("here")
			go GitHub.WillGetClones(connection, &master, db)
		}
	}
	fmt.Println("GitHub Analytics Tracker is Finished")

}

// waitForComplete is intended to run as a rp
func waitForComplete() {
	defer close(finished)
}
