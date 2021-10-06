package config_test

import (
	"os"
	"testing"

	"github.com/nais/aiven-audit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg, err := config.New()
	assert.NoError(t, err)
	assert.Equal(t, &config.AivenAuditConfig{
		AivenAPIToken:   "",
		AuditLogAddress: "audit.nais:6514",
	}, cfg)
}

func TestConfigVariable(t *testing.T) {
	os.Setenv("AIVEN_AUDIT_PAT", "audit pat")
	os.Setenv("AUDIT_LOG_ADDR", "foo:1234")
	cfg, err := config.New()
	assert.NoError(t, err)
	assert.Equal(t, &config.AivenAuditConfig{
		AivenAPIToken:   "audit pat",
		AuditLogAddress: "foo:1234",
	}, cfg)
}

func TestConfigVariableBadAddress(t *testing.T) {
	os.Setenv("AUDIT_LOG_ADDR", "1234")
	cfg, err := config.New()
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "invalid audit log address '1234': expected HOST:PORT")

	os.Setenv("AUDIT_LOG_ADDR", "foo:bar")
	cfg, err = config.New()
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "invalid port number in audit log address 'foo:bar': strconv.Atoi: parsing \"bar\": invalid syntax")

	os.Setenv("AUDIT_LOG_ADDR", "foo:0")
	cfg, err = config.New()
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "invalid port number in audit log address: must be 0 < PORT < 65536")

	os.Setenv("AUDIT_LOG_ADDR", "foo:123456")
	cfg, err = config.New()
	assert.Nil(t, cfg)
	assert.EqualError(t, err, "invalid port number in audit log address: must be 0 < PORT < 65536")
}
