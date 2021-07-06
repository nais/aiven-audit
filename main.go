package main

import (
	"fmt"
	"github.com/aiven/aiven-go-client"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Aiven Audit")
	if err != nil {
		log.Fatal("Could not respond")
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("starting")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func aivenClient() *aiven.Client {
	apiClient, err := aiven.NewTokenClient("test", "")
	if err != nil {
		log.Fatalf("Could not create Aiven Client - err: %s", err)
	}
	return apiClient
}


