package main

import (
	"github.com/rs/zerolog/log"
	"os/exec"
)

func runTests(config Config) {
	for _, impl := range config.Implementations {
		var execFunc func(string, string, ...string) (*exec.Cmd, error)
		switch impl.runtime {
		case Node:
			execFunc = StartNode
		case Python:
			execFunc = StartPython
		case Java:
			execFunc = StartJava
		default:
			log.Error().Str("name", impl.name).Int8("id", int8(impl.runtime)).Msg("Invalid runtime")
		}

		cmd, err := execFunc(impl.path, config.TestsDir+"/"+impl.name)
		if err != nil {
			log.Error().Err(err).Str("name", impl.name).Msg("Starting test command failed")
		}
		cmd.Stdout
	}
}
