package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	securedropbot "github.com/tmc/securedrop-bot"
)

func main() {
	ctx := context.Background()
	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	h, err := securedropbot.NewHandler(ctx)
	if err != nil {
		log.Fatal(err)
	}
	go h.Poll(ctx)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), h))
}
