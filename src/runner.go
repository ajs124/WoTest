package main

import (
	"context"
	"errors"
	"github.com/ajs124/WoTest/clients"
	"github.com/ajs124/WoTest/servers"
	"github.com/rs/zerolog/log"
	"io"
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

func (res *TestResult) waitForOutput(out []io.ReadCloser, match string) (bool, error) {
	rexpr, err := regexp.Compile(match)
	if err != nil {
		log.Err(err)
	}
	buf := make([]byte, 4096)
	err = nil
	m := false
	n := 0
	sout := make([][]byte, 2)
	for err == nil && !m {
		for j, o := range out {
			n, err = o.Read(buf)
			if err != nil {
				log.Debug().Err(err).Msg("Error while trying to read")
				break
			}
			i := 0
			for i < n {
				sout[j] = append(sout[j], buf[i])
				i++
			}
			/* if n > 0 {
				log.Debug().Int("n", n).Str("now", string(sout[j])).Msg("Read bytes")
			} */
			if rexpr.Match(sout[j]) {
				m = true
				break
			}
		}
		time.Sleep(10 * time.Millisecond) // FIXME: why?
	}
	res.stdout += string(sout[0])
	res.stderr += string(sout[1])
	return m, err
}

func runProtocolClientTest(test Test, impl WoTImplementation, config Config, execFunc func(string, string, context.Context, ...string) (*exec.Cmd, []io.ReadCloser, error)) (TestResult, error) {
	result := TestResult{"", "", false}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(test.Timeout))
	defer cancel()
	cmd, out, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, test.Path)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}

	found, err := result.waitForOutput(out, ".* ready.*")
	log.Debug().Bool("found", found).Msg("ready")

	reqUrl, _ := url.Parse(test.ProtocolTestProperties.RequestUrl)

	var client clients.Client

	if test.ProtocolTestProperties.Protocol == ProtoHttp || test.ProtocolTestProperties.Protocol == ProtoHttps {
		client = &clients.HttpClient{}
		r, err := client.Recv(*reqUrl)
		if err != nil {
			log.Err(err).Msg("http client error")
		} else {
			result.succeeded = true
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
			result.succeeded = true
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

func runProtocolServerTest(test Test, impl WoTImplementation, config Config, execFunc func(string, string, context.Context, ...string) (*exec.Cmd, []io.ReadCloser, error)) (TestResult, error) {
	result := TestResult{"", "", false}
	var err error

	if test.ProtocolTestProperties.Protocol == ProtoHttp || test.ProtocolTestProperties.Protocol == ProtoHttps {
		var server servers.HttpServer
		err = server.Listen(test.ProtocolTestProperties.ServeAt)
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
	} else if test.ProtocolTestProperties.Protocol == ProtoMqtt {
		return result, errors.New("not implemented")
	} else {
		// wtf?
		panic(errors.New("test has invalid protocol"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*time.Duration(test.Timeout))
	defer cancel()
	cmd, out, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, test.Path)
	if err != nil {
		log.Error().Err(err).Str("name", impl.Name).Msg("Starting test command failed")
	}

	matches := 0
	for _, mM := range test.ProtocolTestProperties.MustMatch {
		found := false
		found, err = result.waitForOutput(out, mM)
		log.Debug().Bool("found", found).Msg(mM)
		if found {
			matches += 1
		}
	}
	result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
	err = cmd.Process.Kill()
	return result, err
}

func runTest(test Test, impl WoTImplementation, config Config, execFunc func(string, string, context.Context, ...string) (*exec.Cmd, []io.ReadCloser, error)) (TestResult, error) {
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
			fields := make(map[string]interface{})
			if test.Type == TestTypeProtocol {
				fields["protocol"] = test.ProtocolTestProperties.Protocol
				fields["mode"] = test.ProtocolTestProperties.Mode
			}
			logResult.Log().Bool("succeeded", result.succeeded).
				Str("stdout", result.stdout).
				Str("stderr", result.stderr).
				Str("path", test.Path).
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
