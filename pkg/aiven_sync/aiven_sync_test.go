package aiven_sync

import (
	"github.com/aiven/aiven-go-client"
	"testing"
)

type ApiMock struct{}

func (h *ApiMock) GetEventLog(project string) ([]*aiven.ProjectEvent, error) {
	return []*aiven.ProjectEvent{{
		Actor:       "user",
		EventDesc:   "test event",
		EventType:   "test event",
		ServiceName: "test service",
		Time:        "1625742736",
	}}, nil
}

func TestGetEvents(t *testing.T) {
	//api := ApiMock{}
	//sync := NewAivenSync(aiven.ProjectsHandler(api))
}
