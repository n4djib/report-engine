package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadConfigFromFiles(cfg any, envFiles []string) error {
	// disable the env file check for now, we will load the env files in the docker compose file and in the local development environment
	// if len(envFiles) == 0 {
	// 	return errors.New("you have to at east provide one env file")
	// }
	err := LoadAndParseEnv(cfg, envFiles)
	if err != nil {
		return err
	}
	return nil
}

func LoadAndParseEnv(cfg any, envFiles []string) error {
	// load config files to Env system Vars
	// the order of files is important according to if we use load or overload
	// for now it is set to load
	for _, envFile := range envFiles {
		if err := loadEnv(envFile); err != nil {
			fmt.Printf("failed to load env file, error: %v\n", err)
		}
	}
	fmt.Println("Environment variables loaded from config files")

	// parse env variables into the Config struct
	// from the system env variables
	if err := env.Parse(cfg); err != nil {
		// FXME: multiple errors in one string
		return err
	}

	return nil
}

func loadEnv(envFile string) error {
	_, err := os.Stat(envFile)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}
	// Load them into ENV for this process
	// only sets them if they are NOT already set in the system environment
	return godotenv.Load(envFile)
	// this one overloads the env variables even if they are already set in the system environment
	// return godotenv.Overload(envFile)
}
