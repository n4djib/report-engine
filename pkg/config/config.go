package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"

	// "github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func LoadConfigFromFiles(cfg any, envFiles []string) error {
	if len(envFiles) == 0 {
		return errors.New("you have to at east provide one env file")
	}
	err := initAndLoadEnv(cfg, envFiles)
	if err != nil {
		return err
	}
	return nil
}

func loadEnv(envFile string) error {
	_, err := os.Stat(envFile)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatal("File not found\n", err)
	}
	// Load them into ENV for this process
	return godotenv.Load(envFile)
}

func initAndLoadEnv(cfg any, envFiles []string) error {
	// load config files to Env Vars
	for _, envFile := range envFiles {
		if err := loadEnv(envFile); err != nil {
			return err
		}
	}

	fmt.Println("Environment variables loaded from config files")

	// parse env variables into the Config struct
	// from the system? env variables
	if err := env.Parse(cfg); err != nil {
		// multiple errors in one string
		return err
	}

	// debug print pritty the config struct
	// spew.Dump(cfg)

	return nil
}
