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
	//master is a wait group that will wait for all the go routines to finish
	master sync.WaitGroup
	//finished is a channel that will be used to signal the main routine that all the go routines have finished
	finished = make(chan bool, 5)
)

func main() {
	// setup the environment
	stage := Util.MustSetupEnvironment()
	db := Database.MustSetupDB(stage)
	// get all of the repositories from the database and verify the auth tokens for them (check that they work)
	go GithubAuth.ValidateRepositories(db)
	// get all of the repos again to get a count of how many need to be processed
	repos, err := Database.GetAllRepositories(db)
	for err != nil {
		//TODO better error handling
		repos, err = Database.GetAllRepositories(db)
	}
	// add the number of repos to the wait group
	master.Add(len(repos) * 3)
	// start a routine to wait for all of the go routines to finish
	go waitForComplete(&master)

	//named loop to break out of the select statement
L:
	for {
		select {
		case <-finished:
			//break if all of the go routines have finished and signal is sent
			break L
		case connection := <-GithubAuth.Connections:
			//start a separate routine for each type of data
			go GitHub.WillGetClones(connection, &master, db)
			go GitHub.WillGetTrafficReferrers(connection, &master, db)
			go GitHub.WillGetTrafficPaths(connection, &master, db)
		}
	}
	//were done, announce it's finished (graceful chutdown)
	fmt.Println("GitHub Analytics Tracker is Finished")
}

// waitForComplete is intended to run as a routine to wait for all of the go routines to finish
func waitForComplete(group *sync.WaitGroup) {
	group.Wait()
	finished <- true
}
