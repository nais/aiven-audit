package aiven_sync

import (
	"github.com/aiven/aiven-go-client"
	"github.com/nais/aiven-audit/pkg/cef"
	"log"
	"log/syslog"
	"time"
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

func (al AuditLog) Log(events []*aiven.ProjectEvent) ([]*aiven.ProjectEvent, error) {
	auditer, err := syslog.Dial("tcp", al.addr,
		syslog.LOG_INFO|syslog.LOG_DAEMON, al.tag)
	if err != nil {
		return nil, err
	}
	defer func(auditer *syslog.Writer) {
		_ = auditer.Close()
	}(auditer)

	for _, event := range events {
		err = auditer.Info(eventToCef(event).CefString())
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func eventToCef(event *aiven.ProjectEvent) cef.CefRecord {
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
