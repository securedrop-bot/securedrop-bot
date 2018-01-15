package securedropbot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const (
	pollInterval = time.Minute       // TODO: parameterize
	githubOwner  = "securedrop-bot"  // TODO: parameterize?
	githubRepo   = "securedrop-test" // TODO: parameterize?

	// TODO: come up with common structured way to represent thresholds for different policies
	policyNagSubmitterThreshold = 2 * time.Hour
)

func postComment(h *Handler, ctx context.Context, pr *github.PullRequest, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := h.client.Issues.CreateComment(ctx, githubOwner, githubRepo, pr.GetNumber(), comment)
	if err != nil {
		return errors.Wrap(err, "issue posting comment")
	}
	return nil
}

// Handler is the main handler.
type Handler struct {
	logger logrus.FieldLogger
	client *github.Client
}

// NewHandler creates a new Handler.
func NewHandler(ctx context.Context, logger logrus.FieldLogger) (*Handler, error) {
	var tc *http.Client
	if t := os.Getenv("GITHUB_ACCESS_TOKEN"); t != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: t},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	return &Handler{
		logger: logger,
		client: github.NewClient(tc),
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
			h.logger.Infoln("got context done")
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
		h.logger.WithError(err).Warnln("issue listing pull requests")
		return
	}
	fmt.Printf("detected %v pull requests\n", len(prs))
	for _, pr := range prs {
		if err := h.nagSubmitterIfFailed(ctx, pr); err != nil {
			h.logger.WithError(err).WithField("policy", "nagSubmitterIfFailed").Warnln("encountered error")
		}
	}
}

func (h *Handler) nagSubmitterIfFailed(ctx context.Context, pr *github.PullRequest) error {
	logger := h.logger.WithField("policy", "nagSubmitterIfFailed")
	logger.Debugln(pr.GetNumber(), pr.GetState(), pr.GetTitle(), pr.GetBody(), pr.GetUser().GetLogin(), pr.GetStatusesURL())

	statuses, _, err := h.client.Repositories.GetCombinedStatus(ctx, githubOwner, githubRepo, pr.GetHead().GetSHA(), nil)
	if err != nil {
		return errors.Wrap(err, "issue getting status")
	}
	if statuses.GetState() != "failure" {
		logger.WithField("state", statuses.GetState()).Debugln("skipping status")
		return nil
	}
	for _, s := range statuses.Statuses {
		if s.GetState() != "failure" {
			continue
		}
		since := time.Since(s.GetCreatedAt())
		logger.Infoln("created ago:", since)
		if since > policyNagSubmitterThreshold {
			logger.Infoln("would comment if not already")
		}
		body := fmt.Sprintf(`@%v, it looks like there was a test failure, can you please investigate?`, pr.GetUser().GetLogin())

		// TODO: this needs to not post if it's already happened.
		postComment(h, ctx, pr, body)
	}
	return nil
}
