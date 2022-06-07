package main

import (
	"blockchain/pkg/network/response"
	"context"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/reserves/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		response.Send(ctx, 200, struct {
			Message string `json:"message"`
		}{
			Message: "Hello World",
		}, w)
	})
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()
	// start the web server
	log.Println("stop")
}
