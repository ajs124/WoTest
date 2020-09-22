package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ajs124/WoTest/clients"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/url"
	"os/exec"
	"regexp"
	"time"
)

const (
	ModeClient = iota
	ModeServer
)

const (
	ProtoHttp = iota
	ProtoHttps
	ProtoCoap
	ProtoMqtt
)

type TestResult struct {
	stdout    string
	stderr    string
	succeeded bool
}

func waitForOutput(out []io.ReadCloser, match string) (bool, error) {
	cOut := make(chan bool)
	cErr := make(chan bool)
	go func(c chan bool) {
		r, err := regexp.Compile(match)
		buf := make([]byte, 0)
		n, err := out[0].Read(buf)
		r.Match(buf)
	}(cOut)
}

func runTest(test Test, impl WoTImplementation, config Config, execFunc func(string, string, context.Context, ...string) (*exec.Cmd, []io.ReadCloser, error)) (TestResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(test.Timeout))
	defer cancel()
	_, out, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, test.Path)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}
	time.Sleep(time.Second * 1)
	if test.Type == TestTypeProtocol {
		if test.ProtocolTestProperties.Mode == ModeClient {
			reqUrl, _ := url.Parse(test.ProtocolTestProperties.RequestUrl)
			if test.ProtocolTestProperties.Protocol == ProtoHttp || test.ProtocolTestProperties.Protocol == ProtoHttps {
				var client clients.HttpClient
				r, err := client.Recv(*reqUrl)
				if err != nil {
					fmt.Print(err)
				}
				fmt.Print(string(r) + "\n")
			} else if test.ProtocolTestProperties.Protocol == ProtoCoap {
				var client clients.CoapClient
				client.Connect(*reqUrl)
				r, err := client.Recv(*reqUrl)
				if err != nil {
					fmt.Print(err)
				}
				fmt.Print(string(r) + "\n")
			} else if test.ProtocolTestProperties.Protocol == ProtoMqtt {

			} else {
				// wtf?
				panic(errors.New("test has invalid protocol"))
			}
		} else if test.ProtocolTestProperties.Mode == ModeServer {

		} else {
			// wtf?
			panic(errors.New("test has invalid mode"))
		}
	}
	waitForOutput(out)
	cmd.Process.Kill()
	result := TestResult{string(stdout), string(stderr), true}
	return result, nil
}

func runTests(config Config, tests map[string][]Test) {
	for _, impl := range config.Implementations {
		var execFunc func(string, string, context.Context, ...string) (*exec.Cmd, []io.ReadCloser, error)
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
			result, err := runTest(test, impl, config, execFunc)
			if err != nil {
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
