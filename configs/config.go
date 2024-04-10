package configs

import (
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"

	"github.com/elarrg/stori/ledger/internal/adapters/clients/sendgrid"
	"github.com/elarrg/stori/ledger/internal/adapters/db"
)

type Config struct {
	Sendgrid     sendgrid.ClientConfigs `koanf:"sendgrid"`
	PostgresDB   db.PostgresConfig      `koanf:"postgres"`
	Transactions TransactionsConfig     `koanf:"transactions"`
}

type TransactionsConfig struct {
	SourceType   string `koanf:"source-type"`
	SourceFormat string `koanf:"source-format"`
	SourcePath   string `koanf:"source-path"`
}

// Load reads the configs from the available sources, either a YAML formatted file or
// directly from the ENV vars, an ENV config will override a previous one.
func Load() (*Config, error) {
	k := koanf.New(".")

	// load configs from file
	err := k.Load(file.Provider("resources/config.yml"), yaml.Parser())
	if err != nil {
		return nil, err
	}

	// load configs from ENV Vars and merge to file.
	// It replaces an env var from "MY_VAR" to "my.var"
	err = k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(
			strings.ToLower(s),
			"_", ".", -1)
	}), nil)
	if err != nil {
		return nil, err
	}

	// unmarshall configs to struct
	config := new(Config)
	err = k.UnmarshalWithConf("", config, koanf.UnmarshalConf{})
	if err != nil {
		return nil, err
	}

	return config, nil
}
