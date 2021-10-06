package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type AivenAuditConfig struct {
	AivenAPIToken   string `json:"aivenAPIToken"`
	AuditLogAddress string `json:"auditLogAddress"`
}

const (
	AivenAPIToken   = "aivenAPIToken"
	AuditLogAddress = "auditLogAddress"
)

func New() (*AivenAuditConfig, error) {
	viper.SetDefault(AivenAPIToken, "")
	viper.SetDefault(AuditLogAddress, "audit.nais:6514")

	viper.BindEnv(AivenAPIToken, "AIVEN_AUDIT_PAT")
	viper.BindEnv(AuditLogAddress, "AUDIT_LOG_ADDR")

	cfg := &AivenAuditConfig{}

	decoderHook := func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
		dc.ErrorUnused = true
	}

	err := viper.Unmarshal(cfg, decoderHook)
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(cfg.AuditLogAddress, ":")
	if len(tokens) != 2 || len(tokens[0]) == 0 || len(tokens[1]) == 0 {
		return nil, fmt.Errorf("invalid audit log address '%s': expected HOST:PORT", cfg.AuditLogAddress)
	}

	portnumber, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port number in audit log address '%s': %w", cfg.AuditLogAddress, err)
	}

	if portnumber < 1 || portnumber > 65535 {
		return nil, fmt.Errorf("invalid port number in audit log address: must be 0 < PORT < 65536")
	}

	return cfg, err
}
