package aivensync

import (
	"log"
	"log/syslog"
	"time"

	"github.com/nais/aiven-audit/pkg/cef"
)

type AuditLog struct {
	addr string
	tag  string
}

func NewAuditLog(addr, tag string) AuditLog {
	return AuditLog{
		addr,
		tag,
	}
}

func (al AuditLog) Log(event *AivenEvent) error {
	auditer, err := syslog.Dial("tcp", al.addr,
		syslog.LOG_INFO|syslog.LOG_DAEMON, al.tag)

	if err != nil {
		return err
	}

	defer func(auditer *syslog.Writer) {
		_ = auditer.Close()
	}(auditer)

	err = auditer.Info(eventToCef(event).CefString())
	if err != nil {
		return err
	}

	return nil
}

func eventToCef(event *AivenEvent) cef.CefRecord {
	eventTime, err := time.Parse(time.RFC3339, event.Time)
	if err != nil {
		log.Fatalf("Failed to parse event time: %v", err)
	}
	return cef.NewCefRecord(
		0,
		"aiven",
		event.ServiceName,
		"1.0",
		"audit:update",
		event.EventDesc,
		"INFO",
		eventTime.Unix(),
		event.Actor)
}
