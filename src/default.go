package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	configPath := "config.d"
	config, err := loadConfig(configPath)
	if err != nil {
		log.Error().Msgf("Failed to load/parse config file (%s): %s", configPath, err)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))

	tests := make(map[string][]Test)

	for _, i := range config.Implementations {
		log.Debug().Str("Scanning if ImplementationsDir and TestDir exist for implementation", i.Name)
		path := config.ImplementationsDir + "/" + i.Name
		err := checkDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking implementation folder")
		}
		path = config.TestsDir + "/" + i.Name
		err = checkDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking test folder")
		}
		// dhall would have been much nicer, but optional stuff is hard and I'm not good enough at writing it
		t, err := loadTestIndex(path + "/index.json")
		if err != nil {
			log.Error().Err(err).Str("path", path+"/index.d").Msg("Error loading test index")
		}
		tests[i.Name] = t
	}

	runTests(config, tests)
}
