package nais

import (
	"fmt"
	"net/http"
)

func Handlers(mux *http.ServeMux) {
	mux.HandleFunc("/nais/isready", isReady)
	mux.HandleFunc("/nais/isalive", isReady)
}

func isReady(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "OK")
}
