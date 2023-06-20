package Util

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Stage represents the environment the tracker is being run in (e.g. dev, prod, staging etc.)
type Stage string

const (
	Unset Stage = ""
	Dev   Stage = "DEV"
)

var (
	UnsetENVError = errors.New("variable not set")
)

// loadDotENV attempts to load a local development .env file
func loadDotENV() error {
	return godotenv.Load()
}

// GetENV tries to fetch an env variable from particular key and reports the result and success
func GetENV(key string) (string, error) {
	env := os.Getenv(key)
	if env == "" {
		return "", fmt.Errorf("%w: "+key, UnsetENVError)
	}
	return env, nil
}

func GetStage() (Stage, error) {
	val, err := GetENV("STAGE")
	return Stage(val), err
}

// MustSetupEnvironment attempts to set up the environment for the tracker and panics if it is not able to
func MustSetupEnvironment() Stage {
	//TODO make it a switch for different stages
	stage, _ := GetStage()
	if stage == Unset || stage == Dev {
		err := loadDotENV()
		stage, err = GetStage()
		if err != nil {
			log.Panicf("env variable not set: %v", err)
		}
		fmt.Println("GitHub-Analytics-Tracker is running on the DEV(local) environment")
	}
	return stage
	//TODO Production case
}
