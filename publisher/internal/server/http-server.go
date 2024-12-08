package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func HTTPServer(ctx context.Context, queueName string) {
	log.Println("starting HTTP server for health check and metrics")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// should be check to DB for health check
		log.Println("service is healthy")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":      "ok",
			"queue":       queueName,
			"last_update": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error while get health check %v", err)
	}

}
