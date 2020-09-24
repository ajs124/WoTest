package main

import (
	"encoding/json"
	"github.com/philandstuff/dhall-golang/v5"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

const (
	TestTypeProtocol = iota
	TestTypeContents
	TestTypeMeasure
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

type ProtocolTestProperties struct {
	Mode               uint               `json:"mode"`
	Protocol           uint               `json:"protocol"`
	RequestUrl         string             `json:"requestUrl"`
	AuthenticationData AuthenticationData `json:"authentication"`
	ServeAt            string             `json:"serveAt"`
	ServeContent       string             `json:"serveContent"`
	MustMatch          []string           `json:"mustMatch"`
	TlsKey             string             `json:"tlsKey"`
	TlsCert            string             `json:"tlsCert"`
}

type ContentTestProperties struct{}
type MeasureTestProperties struct{}

type Test struct {
	Timeout                uint                   `json:"timeoutSec"`
	Path                   string                 `json:"path"`
	Args                   []string               `json:"args"`
	Type                   uint                   `json:"type"`
	ProtocolTestProperties ProtocolTestProperties `json:"protocolTestProperties"`
	ContentTestProperties  ContentTestProperties  `json:"contentTestProperties"`
	MeasureTestProperties  MeasureTestProperties  `json:"measureTestProperties"`
}

type WoTImplementation struct {
	Name    string  `dhall:"name"`
	Path    string  `dhall:"path"` // this is the "install" or "library" path, not the path of the tests
	Runtime Runtime `dhall:"runtime"`
}

type Config struct {
	TestsDir           string              `dhall:"testsDir"`
	TestResults        string              `dhall:"testResults"`
	ImplementationsDir string              `dhall:"implementationsDir"`
	LogLevel           int8                `dhall:"logLevel"`
	Implementations    []WoTImplementation `dhall:"implementations"`
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

func loadTestIndex(path string) ([]Test, error) {
	tests := make([]Test, 0)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return tests, err
	}
	err = json.Unmarshal(bytes, &tests)
	if err == nil {
		log.Debug().Msgf("Loaded Tests index: %+v", tests)
	}
	return tests, err
}

func checkDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("Error opening")
		return err
	}
	fs, err := f.Stat()
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("Error stat(ing?)")
		return err
	}
	if fs.IsDir() {
		return nil
	}
	return err
}
