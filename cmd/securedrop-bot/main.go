package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	securedropbot "github.com/tmc/securedrop-bot"
)

func main() {
	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	h := securedropbot.NewHandler()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), h))
}
