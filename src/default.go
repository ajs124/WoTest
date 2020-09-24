package main

import (
	. "github.com/ajs124/WoTest/config"
	runner "github.com/ajs124/WoTest/runner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	configPath := "config.d"
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Error().Err(err).Str("filename", configPath).Msg("Failed to load/parse config file")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(zerolog.Level(config.LogLevel))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	os.Rename(config.TestResults, config.TestResults+".old")
	resultFile, err := os.OpenFile(config.TestResults, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		log.Error().Err(err).Str("filename", config.TestResults).Msg("Failed to open results file")
		os.Exit(2)
	}
	resultFile.Truncate(0)
	defer resultFile.Close()
	logResult := zerolog.New(resultFile).With().Logger()

	tests := make(map[string][]Test)

	for _, i := range config.Implementations {
		log.Debug().Str("Scanning if ImplementationsDir and TestDir exist for implementation", i.Name)
		path := config.ImplementationsDir + "/" + i.Name
		err := CheckDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking implementation folder")
		}
		path = config.TestsDir + "/" + i.Name
		err = CheckDir(path)
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Error checking test folder")
		}
		// dhall would have been much nicer, but optional stuff is hard and I'm not good enough at writing it
		t, err := LoadTestIndex(path + "/index.json")
		if err != nil {
			log.Error().Err(err).Str("path", path+"/index.d").Msg("Error loading test index")
		}
		tests[i.Name] = t
	}

	runner.RunTests(config, tests, logResult)
}
