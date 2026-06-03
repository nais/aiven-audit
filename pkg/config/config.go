package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type AivenAuditConfig struct {
	AivenAPIToken string `json:"aivenAPIToken"`
}

const (
	AivenAPIToken = "aivenAPIToken"
)

func New() (*AivenAuditConfig, error) {
	viper.SetDefault(AivenAPIToken, "")
	err := viper.BindEnv(AivenAPIToken, "AIVEN_AUDIT_PAT")
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
