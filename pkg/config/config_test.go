package config_test

import (
	"os"
	"testing"

	"github.com/nais/aiven-audit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	os.Unsetenv("AIVEN_AUDIT_PAT")
	cfg, err := config.New()
	assert.NoError(t, err)
	assert.Equal(t, &config.AivenAuditConfig{
		AivenAPIToken:   "",
	}, cfg)
}

func TestConfigVariable(t *testing.T) {
	os.Setenv("AIVEN_AUDIT_PAT", "audit pat")
	cfg, err := config.New()
	assert.NoError(t, err)
	assert.Equal(t, &config.AivenAuditConfig{
		AivenAPIToken:   "audit pat",
	}, cfg)
}