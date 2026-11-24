// Command server runs the AI Travel Planner API.
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thinhduyng/ai-travel-planner/api/internal/server"
)

func main() {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > 65535 {
			log.Fatal("invalid PORT: must be an integer 1-65535")
		}
		port = p
	}
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           server.New(server.Options{}),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on :%d", port)
	log.Fatal(srv.ListenAndServe())
}
