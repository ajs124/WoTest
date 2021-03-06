package config

import (
	"encoding/json"
	"github.com/philandstuff/dhall-golang/v5"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
)

const (
	TestTypeProtocol = iota
	TestTypeMeasure
	TestTypeContents
)

const (
	ModeClient = iota
	ModeServer
	ModePull
	ModePush
)

const (
	ProtoHttp = iota
	ProtoHttps
	ProtoCoap
	ProtoMqtt
)

// When modifying this, also update config.d
const (
	Node = iota
	Python
	Java
)

type ProtocolTestProperties struct {
	Mode               uint               `json:"mode"`
	Protocol           uint               `json:"protocol"`
	RequestUrl         string             `json:"requestUrl"`
	RequestMustFail    bool               `json:"requestMustFail"`
	AuthenticationData AuthenticationData `json:"authentication"`
	ServeAt            string             `json:"serveAt"`
	ServeContent       string             `json:"serveContent"`
	MustMatch          []string           `json:"mustMatch"`
	TlsKey             string             `json:"tlsKey"`
	TlsCert            string             `json:"tlsCert"`
}

type ContentTestProperties struct{}

type RequestSet struct {
	Num      int `json:"num"`
	Parallel int `json:"parallel"`
}

type MeasureTestProperties struct {
	Protocol           uint               `json:"protocol"`
	Mode               uint               `json:"mode"`
	RequestUrl         string             `json:"requestUrl"`
	AuthenticationData AuthenticationData `json:"authentication"`
	ServeAt            string             `json:"serveAt"`
	ServeContent       string             `json:"serveContent"`
	MustMatch          []string           `json:"mustMatch"`
	TlsKey             string             `json:"tlsKey"`
	TlsCert            string             `json:"tlsCert"`
	RequestSets        []RequestSet       `json:"requestSets"`
}

type Test struct {
	Name                   string                 `json:"name"`
	Timeout                uint                   `json:"timeoutSec"`
	Path                   string                 `json:"path"`
	Args                   []string               `json:"args"`
	Type                   uint                   `json:"type"`
	ProtocolTestProperties ProtocolTestProperties `json:"protocolTestProperties"`
	ContentTestProperties  ContentTestProperties  `json:"contentTestProperties"`
	MeasureTestProperties  MeasureTestProperties  `json:"measureTestProperties"`
}

type WoTImplementation struct {
	Name    string `dhall:"name"`
	Path    string `dhall:"path"` // this is the "install" or "library" path, not the path of the tests
	Runtime uint   `dhall:"runtime"`
}

type Config struct {
	TestsDir           string              `dhall:"testsDir"`
	TestResults        string              `dhall:"testResults"`
	ImplementationsDir string              `dhall:"implementationsDir"`
	LogLevel           int8                `dhall:"logLevel"`
	Implementations    []WoTImplementation `dhall:"implementations"`
}

func LoadConfig(configPath string) (Config, error) {
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

func LoadTestIndex(path string) ([]Test, error) {
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

func CheckDir(path string) error {
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
