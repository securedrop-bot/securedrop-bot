package securedropbot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

const (
	pollInterval = time.Minute       // TODO: parameterize
	githubOwner  = "securedrop-bot"  // TODO: parameterize?
	githubRepo   = "securedrop-test" // TODO: parameterize?
)

// Handler is the main handler.
type Handler struct {
	client *github.Client
}

// NewHandler creates a new Handler.
func NewHandler(ctx context.Context) (*Handler, error) {
	return &Handler{
		client: github.NewClient(nil),
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func (h *Handler) Poll(ctx context.Context) {
	h.poll(ctx)
	t := time.NewTicker(pollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("got context done")
		case <-t.C:
			h.poll(ctx)
		}
	}
}

func (h *Handler) poll(ctx context.Context) {
	log.Println("poll")
	opt := &github.PullRequestListOptions{}
	prs, _, err := h.client.PullRequests.List(ctx, githubOwner, githubRepo, opt)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("detected %v pull requests\n", len(prs))
	for _, pr := range prs {
		fmt.Println(pr)
	}
}
