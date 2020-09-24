package main

import (
	"context"
	"errors"
	"github.com/ajs124/WoTest/clients"
	"github.com/ajs124/WoTest/servers"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
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

func runProtocolClientTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	result := TestResult{"", "", false}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(test.Timeout))
	defer cancel()
	stdout := make([]byte, 0)
	stderr := make([]byte, 0)
	cmd, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, &stdout, &stderr, test.Path)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}

	matches := 0
	for _, mM := range test.ProtocolTestProperties.MustMatch {
		found := result.waitForOutput(ctx, &stdout, &stderr, mM)
		log.Debug().Str("mM", mM).Bool("found", found).Msg("found")
		if found {
			matches += 1
		}
	}
	log.Debug().Int("matches", matches).Int("out of", len(test.ProtocolTestProperties.MustMatch)).Msg("")
	reqUrl, _ := url.Parse(test.ProtocolTestProperties.RequestUrl)

	var client clients.Client

	if test.ProtocolTestProperties.Protocol == ProtoHttp || test.ProtocolTestProperties.Protocol == ProtoHttps {
		client = &clients.HttpClient{}
		r, err := client.Recv(*reqUrl)
		if err != nil {
			log.Err(err).Msg("http client error")
			time.Sleep(time.Minute)
		} else {
			result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
			log.Debug().Str("http client output", string(r))
		}
	} else if test.ProtocolTestProperties.Protocol == ProtoCoap {
		client = &clients.CoapClient{}
		_, err = client.Connect(*reqUrl)
		if err != nil {
			log.Err(err).Msg("coap client failed to connect")
		}
		r, err := client.Recv(*reqUrl)
		if err != nil {
			log.Err(err).Msg("coap client failed to GET data")
		} else {
			result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
			log.Debug().Str("coap client output", string(r))
		}
	} else if test.ProtocolTestProperties.Protocol == ProtoMqtt {
		return result, errors.New("not implemented")
	} else {
		// wtf?
		panic(errors.New("test has invalid protocol"))
	}

	err = cmd.Process.Kill()
	return result, nil
}

func runProtocolServerTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	result := TestResult{"", "", false}
	var err error

	if test.ProtocolTestProperties.Protocol == ProtoHttp || test.ProtocolTestProperties.Protocol == ProtoHttps {
		var server servers.HttpServer
		if test.ProtocolTestProperties.Protocol == ProtoHttp {
			err = server.Listen(test.ProtocolTestProperties.ServeAt, false, "", "")
		} else { // TLS
			err = server.Listen(test.ProtocolTestProperties.ServeAt, true, test.ProtocolTestProperties.TlsCert, test.ProtocolTestProperties.TlsKey)
		}
		defer server.Stop()
		if err != nil {
			log.Err(err).Msg("http server error")
		}

		server.SetHandleFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: add support for other types of content
			http.ServeFile(w, r, test.ProtocolTestProperties.ServeContent)
		})
	} else if test.ProtocolTestProperties.Protocol == ProtoCoap {
		var server servers.CoapServer
		err = server.Listen(test.ProtocolTestProperties.ServeAt)
		defer server.Stop()
		if err != nil {
			log.Err(err).Msg("coap server error")
		}
		server.SetHandleFunc(func(w mux.ResponseWriter, r *mux.Message) {
			server.ServeFile(w, r, test.ProtocolTestProperties.ServeContent)
		})
	} else if test.ProtocolTestProperties.Protocol == ProtoMqtt {
		return result, errors.New("not implemented")
	} else {
		// wtf?
		panic(errors.New("test has invalid protocol"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(test.Timeout))
	defer cancel()
	stdout := make([]byte, 0)
	stderr := make([]byte, 0)
	cmd, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, &stdout, &stderr, test.Path)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}

	matches := 0
	for _, mM := range test.ProtocolTestProperties.MustMatch {
		found := result.waitForOutput(ctx, &stdout, &stderr, mM)
		log.Debug().Str("mM", mM).Bool("found", found).Msg("found")
		if found {
			matches += 1
		}
	}
	result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
	err = cmd.Process.Kill()
	return result, err
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

func runTests(config Config, tests map[string][]Test) {
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
