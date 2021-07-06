package api_client

import (
	"github.com/aiven/aiven-go-client"
	"log"
)

func CreateAivenClient() *aiven.Client {
	apiClient, err := aiven.NewTokenClient("test", "")
	if err != nil {
		log.Fatalf("Could not create Aiven Client - err: %s", err)
	}
	return apiClient
}
