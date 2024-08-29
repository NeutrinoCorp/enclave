package transport

import (
	"github.com/caarlos0/env/v11"
	"github.com/neutrinocorp/nolan/collection/set"
)

type ConfigHTTP struct {
	Address                 string   `env:"HTTP_SERVER_ADDRESS" envDefault:":8080"`
	AuthenticationWhitelist []string `env:"HTTP_SERVER_AUTHENTICATION_WHITELIST" envDefault:"/healthz,/readiness"`

	AuthenticationWhitelistSet set.Set[string]
}

func NewConfigHTTP() (ConfigHTTP, error) {
	cfg, err := env.ParseAs[ConfigHTTP]()
	if err != nil {
		return ConfigHTTP{}, err
	}

	cfg.AuthenticationWhitelistSet = set.HashSet[string]{}
	for _, item := range cfg.AuthenticationWhitelist {
		cfg.AuthenticationWhitelistSet.Add(item)
	}
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
