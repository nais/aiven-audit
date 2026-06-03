package aivensync

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type aivenEventsResponse struct {
	Events []*AivenEvent `json:"events"`
}

type AivenEvent struct {
	ID          string `json:"id"`
	Actor       string `json:"actor"`
	EventDesc   string `json:"event_desc"`
	EventType   string `json:"event_type"`
	ServiceName string `json:"service_name"`
	Time        string `json:"time"`
}

func (a *AivenEvent) Equals(b *AivenEvent) bool {
	return a.ID == b.ID
}

type AivenProject struct {
	ProjectName string `json:"project_name"`
}

type aivenProjectsResponse struct {
	Projects []*AivenProject `json:"projects"`
}

type AivenTransport struct {
	aivenToken string
}

func (at AivenTransport) Client() *http.Client {
	return &http.Client{Transport: at}
}

func (at AivenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+at.aivenToken)
	return http.DefaultTransport.RoundTrip(req)
}

type AivenClient interface {
	GetProjectEvents(project string) ([]*AivenEvent, error)
	GetProjects() ([]*AivenProject, error)
}

type aivenClient struct {
	client *http.Client
}

func closeResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		slog.Warn("Unable to close response body", "error", err)
	}
}

func (a *aivenClient) GetProjects() ([]*AivenProject, error) {
	response, err := a.client.Get("https://api.aiven.io/v1/project")
	if err != nil {
		return nil, err
	}

	defer closeResponseBody(response.Body)

	err = checkForHttpError(response)
	if err != nil {
		return nil, err
	}

	projects := &aivenProjectsResponse{}
	err = json.NewDecoder(response.Body).Decode(projects)
	if err != nil {
		return nil, err
	}

	return projects.Projects, nil
}

func (a *aivenClient) GetProjectEvents(project string) ([]*AivenEvent, error) {
	response, err := a.client.Get(fmt.Sprintf("https://api.aiven.io/v1/project/%s/events?limit=1000", project))
	if err != nil {
		return nil, err
	}

	defer closeResponseBody(response.Body)

	err = checkForHttpError(response)
	if err != nil {
		return nil, err
	}

	events := &aivenEventsResponse{}
	err = json.NewDecoder(response.Body).Decode(events)
	if err != nil {
		return nil, err
	}

	return events.Events, nil
}

func NewAivenClient(token string) AivenClient {
	client := &http.Client{
		Transport: AivenTransport{
			aivenToken: token,
		},
	}

	return &aivenClient{client: client}
}

func checkForHttpError(response *http.Response) error {
	if response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusOK || response.StatusCode == http.StatusAccepted {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("parse http error response body: %s", err)
	}

	return fmt.Errorf("http response: %d, message: %s", response.StatusCode, string(body))
}
