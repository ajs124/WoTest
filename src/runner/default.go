package runner

import (
	"context"
	"errors"
	. "github.com/ajs124/WoTest/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os/exec"
	"regexp"
	"time"
)

type TestResult struct {
	stdout    string
	stderr    string
	succeeded bool
}

type ExecFunc func(string, string, context.Context, *[]byte, *[]byte, ...string) (*exec.Cmd, error)

func (res *TestResult) waitForOutput(ctx context.Context, stdout, stderr *[]byte, match string) bool {
	rexpr, err := regexp.Compile(match)
	if err != nil {
		log.Err(err)
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C: // maybe the regex will match now?
			res.stdout = string(*stdout)
			res.stderr = string(*stderr)
			if rexpr.Match(*stdout) || rexpr.Match(*stderr) {
				return true
			}
		}
	}
}

func runTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	if test.Type == TestTypeProtocol {
		if test.ProtocolTestProperties.Mode == ModeClient {
			return runProtocolClientTest(test, impl, config, execFunc)
		} else if test.ProtocolTestProperties.Mode == ModeServer {
			return runProtocolServerTest(test, impl, config, execFunc)
		} else {
			// wtf?
			panic(errors.New("test has invalid mode"))
		}
	}
	panic(errors.New("cannot run this type of test"))
}

func RunTests(config Config, tests map[string][]Test, logResult zerolog.Logger) {
	for _, impl := range config.Implementations {
		var execFunc ExecFunc
		switch impl.Runtime {
		case Node:
			execFunc = StartNode
		case Python:
			execFunc = StartPython
		case Java:
			execFunc = StartJava
		default:
			log.Error().Str("name", impl.Name).Int8("id", int8(impl.Runtime)).Msg("Invalid runtime")
		}

		for _, test := range tests[impl.Name] {
			log.Info().Str("implementation", impl.Name).Str("path", test.Path).Msg("Starting test")
			result, err := runTest(test, impl, config, execFunc)
			fields := make(map[string]interface{})
			if test.Type == TestTypeProtocol {
				fields["protocol"] = test.ProtocolTestProperties.Protocol
				fields["mode"] = test.ProtocolTestProperties.Mode
			}
			logResult.Log().Bool("succeeded", result.succeeded).
				Str("stdout", result.stdout).
				Str("stderr", result.stderr).
				Str("path", test.Path).
				Strs("args", test.Args).
				Str("implementation", impl.Name).
				Fields(fields).
				Uint("type", test.Type).Msg("")
			if err != nil || !result.succeeded {
				log.Error().
					Err(err).
					Str("path", test.Path).
					Msg("Test failed")
			} else {
				log.Info().
					Bool("succeeded", result.succeeded).
					Str("stdout", result.stdout).
					Str("stderr", result.stderr).
					Str("path", test.Path).
					Uint("type", test.Type).
					Msg("Ran test")
			}
		}
	}
}
