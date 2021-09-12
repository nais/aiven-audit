package config

import "os"

type AivenAuditConfig struct {
	AivenAPIToken string
	AuditLogAddr  string
}

func ConfigFromEnv() AivenAuditConfig {
	aivenPat := os.Getenv("AIVEN_AUDIT_PAT")
	auditAddr := os.Getenv("AUDIT_LOG_ADDR")
	return AivenAuditConfig{
		aivenPat,
		auditAddr,
	}
}
