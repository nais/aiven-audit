package config

import (
	"os"
)

type AivenAuditConfig struct {
	AivenAPIToken string
	AuditLogAddr  string
}

func FromEnv() AivenAuditConfig {
	aivenPat := os.Getenv("AIVEN_AUDIT_PAT")
	auditAddr := os.Getenv("AUDIT_LOG_ADDR")
	if auditAddr == "" {
		auditAddr = "audit.nais"
	}

	return AivenAuditConfig{
		aivenPat,
		auditAddr,
	}
}
