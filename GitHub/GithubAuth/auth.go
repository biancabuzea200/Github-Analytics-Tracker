package GithubAuth

import (
	"GitHub-Analytics-Tracker/Database"
	"GitHub-Analytics-Tracker/Util"
	"context"
	"fmt"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"sync"
)

var (
	// Connections is a channel where valid connections are placed
	Connections = make(chan *Connection, 50)
	//TODO make buffer dynamic?
)

type Connection struct {
	Repository *Database.Repository
	Id         string
	Client     *github.Client
	//TODO implement some rate limiting logic if necessary
	//https://docs.github.com/en/rest/rate-limit?apiVersion=2022-11-28
}

// getToken gets the api (personal access Connection) for the appropriate Repository
func mustGetTokenENV(id string) string {
	// PAT stands for personal access connection (https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens)
	PAT, err := Util.GetENV(id)
	if err != nil {
		panic(err)
	}
	return PAT
}

// ValidateRepositories is designed to run as a routine that generates valid GitHub clients and sends them to Connections when they are ready
func ValidateRepositories(db *gorm.DB) {
	fmt.Println("Creating connections for all repositories")
	repositories, err := Database.GetAllRepositories(db)
	//TODO better error handling here
	for err != nil {
		repositories, err = Database.GetAllRepositories(db)
	}
	var wg sync.WaitGroup
	wg.Add(len(repositories))
	for _, repo := range repositories {
		go mustCreateValidConnection(repo, &wg)
	}
	wg.Wait()
	fmt.Println("Finished generating connections for all repositories")
}

// mustCreateValidConnection creates a connection that can be used to interact with the GitHub API
func mustCreateValidConnection(repository *Database.Repository, wg *sync.WaitGroup) {
	defer wg.Done()
	PAT := mustGetTokenENV(repository.TokenKey)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: PAT},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	err := validateClient(client)
	if err != nil {
		panic(err)
	}

	Connections <- &Connection{
		Id:         repository.TokenKey,
		Client:     client,
		Repository: repository,
	}
}

func validateClient(client *github.Client) error {
	_, _, err := client.Repositories.GetByID(context.Background(), 1)
	return err
}
