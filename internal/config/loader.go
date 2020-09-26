package config

import (
	"fmt"
	"io"

	"github.com/spf13/viper"
)

var applyEnvOverrides = envOverride(envOverrides)

// Load loads config from confFile into Config. The file must contain valid YAML.
func Load(confFile io.Reader) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// set env overrides so that secrets can be passed-in as env variables
	applyEnvOverrides(v)

	if err := v.ReadConfig(confFile); err != nil {
		return nil, fmt.Errorf("cannot read config: %v", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal to config struct: %v", err)
	}

	return &c, nil
}

func envOverride(overrides map[string]string) func(*viper.Viper) {
	return func(v *viper.Viper) {
		for key, override := range overrides {
			// error is returned only when len(input) == 0, which can be safely ignored
			_ = v.BindEnv(key, override)
		}
	}
}
