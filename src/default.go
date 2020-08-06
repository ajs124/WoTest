package main

import (
	"github.com/philandstuff/dhall-golang/v4"
	"io/ioutil"
	"log"
	"os"
)

type WoTConsumer interface {
}

type WoTProducer interface {
}

type WoTImplementationInterface interface {
	initialize() error
	consumer() (WoTConsumer, error)
	producer() (WoTProducer, error)
}

type WoTImplementation struct {
	name string
}

type Config struct {
	TestsDir string
}

func loadConfig(configPath string) (Config, error) {
	var config Config
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = dhall.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}
	return config, err
}

func main() {
	configPath := "config.d"
	config, err := loadConfig(configPath)
	if err != nil {
		log.Printf("Failed to load/parse config file (%s): %s", configPath, err)
		os.Exit(1)
	}

	var WoTImplementations []WoTImplementation

	testsDir := config.TestsDir
	folders, err := ioutil.ReadDir(testsDir)
	if err != nil {
		log.Printf("Failed to read tests folder (%s): %s", testsDir, err)
		os.Exit(2)
	}
	for _, d := range folders {
		if !d.IsDir() {
			log.Printf("Non-directory file in tests folder (%s): %s", testsDir, d.Name())
		} else {
			log.Printf("Adding implementation from tests folder (%s): %s", testsDir, d.Name())
			WoTImplementations = append(WoTImplementations, WoTImplementation{name: d.Name()})
		}
	}
}
