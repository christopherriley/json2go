package main

import (
	"fmt"
	"net/http"

	"github.com/christopherriley/json2go/api"
)

func main() {
	mux := http.NewServeMux()
	apiMux := api.NewApi().GetServeMux()

	mux.Handle("/", apiMux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
