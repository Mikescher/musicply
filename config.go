package mply

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gogs.mikescher.com/BlackForestBytes/goext/confext"
	"os"
	"time"
)

const APILevel = 1

type Config struct {
	Namespace       string
	ReturnRawErrors bool          `env:"RETURN_RAW_ERRORS"`
	GinDebug        bool          `env:"GIN_DEBUG"`
	Custom404       bool          `env:"CUSTOM_404"`
	LogLevel        zerolog.Level `env:"LOGLEVEL"`
	ServerIP        string        `env:"SERVER_IP"`
	ServerPort      string        `env:"SERVER_PORT"`
	RequestTimeout  time.Duration `env:"REQUEST_TIMEOUT"`
	Cors            bool          `env:"CORS"`
	LiveReload      *string       `env:"LIVE_RELOAD"`
}

var defaultConf = Config{
	Namespace:       "default",
	ReturnRawErrors: false,
	GinDebug:        false,
	Custom404:       false,
	ServerIP:        "0.0.0.0",
	ServerPort:      "8000",
	RequestTimeout:  16 * time.Second,
	LogLevel:        zerolog.WarnLevel,
	Cors:            true,
	LiveReload:      nil,
}

var localConf = Config{
	Namespace:       "local",
	ReturnRawErrors: true,
	GinDebug:        true,
	Custom404:       false,
	ServerIP:        "0.0.0.0",
	ServerPort:      "8000",
	RequestTimeout:  60 * time.Second,
	LogLevel:        zerolog.DebugLevel,
	Cors:            true,
	LiveReload:      nil,
}

var allConfigs = map[string]Config{
	"default": defaultConf,
	"local":   localConf,
}

var Conf Config

func getConfig(ns string) (Config, bool) {
	if ns == "" {
		ns = "default"
	}
	if conf, ok := allConfigs[ns]; ok {
		err := confext.ApplyEnvOverrides("", &conf, "_")
		if err != nil {
			panic(err)
		}
		return conf, true
	}
	return Config{}, false //nolint:exhaustruct
}

func init() {
	ns := os.Getenv("CONF_NS")
	conf, ok := getConfig(ns)
	if !ok {
		log.Fatal().Str("ns", ns).Msg("Unknown config-namespace")
	}
	Conf = conf
}
