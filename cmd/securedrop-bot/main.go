package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	securedropbot "github.com/tmc/securedrop-bot"
)

var flagVerbose = flag.Bool("v", false, "be verbose")

func main() {
	flag.Parse()
	logger := logrus.New()
	if *flagVerbose {
		logger.SetLevel(logrus.DebugLevel)
		logrus.SetLevel(logrus.DebugLevel)
	}
	ctx := context.Background()
	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	h, err := securedropbot.NewHandler(ctx, logger)
	if err != nil {
		log.Fatal(err)
	}
	go h.Poll(ctx)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), h))
}
