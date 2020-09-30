package runner

import (
	"context"
	"errors"
	"github.com/ajs124/WoTest/clients"
	. "github.com/ajs124/WoTest/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type TestMeasurements struct {
	total  time.Duration
	median time.Duration
	nfth   time.Duration
	nnth   time.Duration
	mean   time.Duration
	stddev time.Duration
	raw    []clients.MeasurementData
}

func (tm *TestMeasurements) MarshalJSON() ([]byte, error) {
	output := "{"
	output += "\"total\":\"" + tm.total.String() + "\","
	output += "\"median\":\"" + tm.median.String() + "\","
	output += "\"95%\":\"" + tm.nfth.String() + "\","
	output += "\"99%\":\"" + tm.nnth.String() + "\","
	output += "\"mean\":\"" + tm.mean.String() + "\","
	output += "\"stddev\":\"" + tm.stddev.String() + "\","
	output += "\"raw\":["
	for i, r := range tm.raw {
		output += "{ \"start\": \"" + r.StartTime.String() + "\","
		output += "\"stop\": \"" + r.StopTime.String() + "\","
		output += "\"duration\": \"" + strconv.FormatInt(r.StopTime.Sub(r.StartTime).Microseconds(), 10) + "\","
		output += "\"size\": \"" + strconv.FormatUint(uint64(r.Size), 10) + "\"}"
		if i != len(tm.raw)-1 {
			output += ","
		}
	}
	output += "]"
	output += "}"
	return []byte(output), nil
}

type TestResult struct {
	name         string
	stdout       string
	stderr       string
	succeeded    bool
	measurements map[int]*TestMeasurements
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
	} else if test.Type == TestTypeMeasure {
		return runMeasureTest(test, impl, config, execFunc)
	} else if test.Type == TestTypeContents {
		panic(errors.New("contents tests are not yet implemented"))
	} else {
		panic(errors.New("cannot run this type of test. invalid test type?"))
	}
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
			fields["measurements"] = result.measurements
			if test.Name == "" {
				test.Name = impl.Name + " " + test.Path
			}
			if test.Type == TestTypeProtocol {
				fields["protocol"] = test.ProtocolTestProperties.Protocol
				fields["mode"] = test.ProtocolTestProperties.Mode
				fields["requestMustFail"] = test.ProtocolTestProperties.RequestMustFail
			}
			logResult.Log().Bool("succeeded", result.succeeded).
				Str("stdout", result.stdout).
				Str("stderr", result.stderr).
				Str("path", test.Path).
				Strs("args", test.Args).
				Str("implementation", impl.Name).
				Str("name", test.Name).
				Fields(fields).
				Uint("type", test.Type).Msg("")
			if err != nil || !result.succeeded {
				log.Error().
					Err(err).
					Str("path", test.Path).
					Str("name", test.Name).
					Msg("Test failed")
			} else {
				log.Info().
					Bool("succeeded", result.succeeded).
					Str("stdout", result.stdout).
					Str("stderr", result.stderr).
					Str("path", test.Path).
					Str("name", test.Name).
					Uint("type", test.Type).
					Fields(fields).
					Msg("Ran test")
			}
		}
	}
}
