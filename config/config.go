package config

import (
	"github.com/gin-gonic/gin"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	*koanf.Koanf
}

func NewConfig(fname string) *Config {
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var conf = koanf.New(".")

	conf.Load(file.Provider(fname), yaml.Parser())

	if conf.String("profile") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Config{conf}
}
