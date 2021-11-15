package main

import (
	"flag"
	"fmt"
	"net/http"
	"storage/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	port := flag.Int("port", 8080, "http port")
	flag.Parse()
	api := server.NewApi()
	r := api.NewRouter()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), r); err != nil {
		log.Fatalf("failed to serve api on port %d: %v", port, err)
	}
}
