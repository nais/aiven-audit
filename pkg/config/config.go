package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type AivenAuditConfig struct {
	AivenAPIToken string `json:"aivenAPIToken"`
	Tenant        string `json:"tenant"`
}

const (
	AivenAPIToken = "aivenAPIToken"
	Tenant        = "tenant"
)

func New() (*AivenAuditConfig, error) {
	viper.SetDefault(AivenAPIToken, "")
	viper.SetDefault(Tenant, "")
	err := viper.BindEnv(AivenAPIToken, "AIVEN_AUDIT_PAT")
	err = viper.BindEnv(Tenant, "AIVEN_TENANT")

	if err != nil {
		return nil, err
	}

	cfg := &AivenAuditConfig{}

	decoderHook := func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
		dc.ErrorUnused = true
	}

	err = viper.Unmarshal(cfg, decoderHook)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
