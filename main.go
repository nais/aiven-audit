package main

import (
	"fmt"
	api_client "github.com/nais/aiven-audit/pkg/api-client"
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
	_ = api_client.CreateAivenClient()
	log.Fatal(http.ListenAndServe(":8080", nil))
}



