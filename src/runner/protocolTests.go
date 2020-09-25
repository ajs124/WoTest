package runner

import (
	"context"
	"errors"
	"github.com/ajs124/WoTest/clients"
	. "github.com/ajs124/WoTest/config"
	"github.com/ajs124/WoTest/servers"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"time"
)

func runProtocolClientTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	var result TestResult
	result.measurements = make(map[int]*TestMeasurements)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(test.Timeout))
	defer cancel()
	stdout := make([]byte, 0)
	stderr := make([]byte, 0)
	cmd, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, &stdout, &stderr, append([]string{test.Path}, test.Args...)...)
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
		r, md, err := client.Recv(*reqUrl)
		if err != nil {
			log.Err(err).Msg("http client error")
		} else {
			result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
			m := TestMeasurements{}
			m.raw = []clients.MeasurementData{md}
			result.measurements[0] = &m
			log.Debug().Str("http client output", string(r))
		}
	} else if test.ProtocolTestProperties.Protocol == ProtoCoap {
		client = &clients.CoapClient{}
		_, err = client.Connect(*reqUrl)
		if err != nil {
			log.Err(err).Msg("coap client failed to connect")
		}
		r, md, err := client.Recv(*reqUrl)
		if err != nil && !test.ProtocolTestProperties.RequestMustFail {
			log.Err(err).Msg("coap client failed to GET data")
		} else {
			result.succeeded = len(test.ProtocolTestProperties.MustMatch) == matches
			m := TestMeasurements{}
			m.raw = []clients.MeasurementData{md}
			result.measurements[0] = &m
			log.Debug().Str("coap client output", string(r))
		}
	} else if test.ProtocolTestProperties.Protocol == ProtoMqtt {
		return result, errors.New("not implemented")
	} else {
		// wtf?
		panic(errors.New("test has invalid protocol"))
	}

	result.stdout = string(stdout)
	result.stderr = string(stderr)
	err = cmd.Process.Kill()
	return result, nil
}

func runProtocolServerTest(test Test, impl WoTImplementation, config Config, execFunc ExecFunc) (TestResult, error) {
	var result TestResult
	result.measurements = make(map[int]*TestMeasurements)
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
	cmd, err := execFunc(impl.Path, config.TestsDir+"/"+impl.Name, ctx, &stdout, &stderr, append([]string{test.Path}, test.Args...)...)
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

	result.stdout = string(stdout)
	result.stderr = string(stderr)
	err = cmd.Process.Kill()
	return result, err
}
