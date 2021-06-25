package config

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	*koanf.Koanf
}

func NewConfig(fname string, envPrefix string) *Config {
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var conf = koanf.New(".")

	// Load JSON config.
	if err := conf.Load(file.Provider(fname), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	conf.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)

	if conf.String("profile") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Config{conf}
}
