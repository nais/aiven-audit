package nais

import (
	"fmt"
	"net/http"
)

func InitNaisHandlers() {
	http.HandleFunc("/nais/isready", isReady)
	http.HandleFunc("/nais/isalive", isReady)
}

func isReady(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
