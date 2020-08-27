package main

import (
	"github.com/philandstuff/dhall-golang/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

type WoTImplementation struct {
	name    string
	path    string // this is the "install" or "library" path, not the path of the
	runtime Runtime
}

type Config struct {
	TestsDir           string
	ImplementationsDir string
	LogLevel           int8
	Implementations    []WoTImplementation
}

func loadConfig(configPath string) (Config, error) {
	var config Config
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = dhall.Unmarshal(bytes, &config)
	if err == nil {
		log.Debug().Msgf("Loaded Config: %+v", config)
	}
	return config, err
}

func checkDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("Error opening")
	}
	fs, err := f.Stat()
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("Error stat(ing?)")
	}
	if fs.IsDir() {
		return nil
	}
	return err
}

func main() {
	configPath := "config.d"
	config, err := loadConfig(configPath)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))
		log.Error().Msgf("Failed to load/parse config file (%s): %s", configPath, err)
		os.Exit(1)
	}

	for _, i := range config.Implementations {
		log.Debug().Str("Scanning if ImplementationsDir and TestDir exist for implementation", i.name)
		path := config.ImplementationsDir + "/" + i.name
		err := checkDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking implementation folder")
		}
		path = config.TestsDir + "/" + i.name
		err = checkDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking test folder")
		}
	}

	runTests(config)
}
