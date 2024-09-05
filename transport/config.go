package transport

import (
	"github.com/caarlos0/env/v11"
	"github.com/emirpasic/gods/v2/sets"
	"github.com/emirpasic/gods/v2/sets/hashset"
)

type ConfigHTTP struct {
	Address                 string   `env:"HTTP_SERVER_ADDRESS" envDefault:":8080"`
	AuthenticationWhitelist []string `env:"HTTP_SERVER_AUTHENTICATION_WHITELIST" envDefault:"/healthz,/readiness"`

	AuthenticationWhitelistSet sets.Set[string]
}

func NewConfigHTTP() (ConfigHTTP, error) {
	cfg, err := env.ParseAs[ConfigHTTP]()
	if err != nil {
		return ConfigHTTP{}, err
	}

	cfg.AuthenticationWhitelistSet = hashset.New(cfg.AuthenticationWhitelist...)
	return cfg, nil
}

type ConfigActuatorHTTP struct {
	ActuatorRoleAllowlist []string `env:"HTTP_SERVER_ACTUATOR_ROLE_ALLOWLIST"`
}

func NewConfigActuatorHTTP() (ConfigActuatorHTTP, error) {
	cfg, err := env.ParseAs[ConfigActuatorHTTP]()
	if err != nil {
		return ConfigActuatorHTTP{}, err
	}
	return cfg, nil
}
