package config

import "os"
import "strings"

type AivenAuditConfig struct {
	AivenAPIToken string
	AuditLogAddr  string
	DbHost        string
	Projects      []string
}

func ConfigFromEnv() AivenAuditConfig {
	aivenPat := os.Getenv("AIVEN_AUDIT_PAT")
	auditAddr := os.Getenv("AUDIT_LOG_ADDR")
	dbHost := os.Getenv("NAIS_DATABASE_AIVENAUDIT_EVENTHASHDB_URL")
	strProjects := os.Getenv("AIVEN_PROJECTS")
	projects := strings.Split(strProjects, " ")
	return AivenAuditConfig{
		aivenPat,
		auditAddr,
		dbHost,
		projects,
	}
}
