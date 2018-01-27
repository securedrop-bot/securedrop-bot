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

var (
	githubOwner = getEnv("SDBOT_REPO_OWNER", "securedrop-bot")
	githubRepo  = getEnv("SDBOT_REPO_NAME", "securedrop-test")
	botUsername = getEnv("SDBOT_USERNAME", "securedrop-bot")
)

const (
	pollInterval = time.Minute // TODO: parameterize

	// TODO: come up with common structured way to represent thresholds for different policies
	policyNagSubmitterThreshold               = 2 * time.Hour
	policyNagReviewerThreshold                = 48 * time.Hour
	policyNagSubmitterReviewCommentsThreshold = 48 * time.Hour
	policyNagMaintainerToMergeThreshold       = 12 * time.Hour
)

// Handler is the main handler.
type Handler struct {
	logger logrus.FieldLogger
	client githubClient
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
		client: &realGithubClient{github.NewClient(tc)},
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
	prs, _, err := h.client.GetPullRequestsService().List(ctx, githubOwner, githubRepo, opt)
	if err != nil {
		h.logger.WithError(err).Warnln("issue listing pull requests")
		return
	}
	fmt.Printf("detected %v pull requests\n", len(prs))
	for _, pr := range prs {
		if err := h.nagSubmitterIfFailed(ctx, pr); err != nil {
			h.logger.WithError(err).WithField("policy", "nagSubmitterIfFailed").Warnln("encountered error")
		}
		if err := h.nagReviewerIfSlow(ctx, pr); err != nil {
			h.logger.WithError(err).WithField("policy", "nagReviewerIfSlow").Warnln("encountered error")
		}
		if err := h.nagMaintainerForMerge(ctx, pr); err != nil {
			h.logger.WithError(err).WithField("policy", "nagMaintainerForMerge").Warnln("encountered error")
		}
	}
}

func (h *Handler) nagSubmitterIfFailed(ctx context.Context, pr *github.PullRequest) error {
	logger := h.logger.WithField("policy", "nagSubmitterIfFailed")
	logger.Debugln(pr.GetNumber(), pr.GetState(), pr.GetTitle(), pr.GetBody(), pr.GetUser().GetLogin(), pr.GetStatusesURL())

	statuses, _, err := h.client.GetRepositoriesService().GetCombinedStatus(ctx, githubOwner, githubRepo, pr.GetHead().GetSHA(), nil)
	if err != nil {
		return errors.Wrap(err, "issue getting status")
	}
	if statuses.GetState() != "failure" {
		logger.WithField("state", statuses.GetState()).Debugln("skipping status")
		return nil
	}
	var sinceLastFailure time.Duration
	for _, s := range statuses.Statuses {
		if s.GetState() != "failure" {
			continue
		}
		sinceLastFailure = time.Since(s.GetCreatedAt())
		break
	}
	comments, err := h.getComments(ctx, pr)
	if err != nil {
		return errors.Wrap(err, "issue getting PR comments")
	}

	lastTimeBotCommented := time.Time{}
	for _, comment := range comments {
		commentUser := *comment.User.Login
		if commentUser == botUsername {
			lastTimeBotCommented = *comment.CreatedAt
		}
	}
	if sinceLastFailure > policyNagSubmitterThreshold && lastTimeBotCommented.IsZero() {
		body := fmt.Sprintf(`@%v, it looks like there was a test failure, can you please investigate?`, pr.GetUser().GetLogin())

		return h.postComment(ctx, pr, body)
	}
	return nil
}

func (h *Handler) nagMaintainerForMerge(ctx context.Context, pr *github.PullRequest) error {
	reviews, _ := h.getReviews(ctx, pr)

	prWasApproved := false
	lastReviewer := ""
	for _, review := range reviews {
		if *review.State != "APPROVED" {
			return nil
		}
		prWasApproved = true
		lastReviewer = *review.User.Login
	}

	comments, err := h.getComments(ctx, pr)
	if err != nil {
		return errors.Wrap(err, "issue getting PR comments")
	}

	lastTimeBotCommented := time.Time{}
	for _, comment := range comments {
		commentUser := *comment.User.Login
		if commentUser == botUsername {
			lastTimeBotCommented = *comment.CreatedAt
		}
	}

	body := fmt.Sprintf("@%v, can we merge this PR?", lastReviewer)

	if prWasApproved && time.Since(lastTimeBotCommented) > policyNagMaintainerToMergeThreshold {
		return h.postComment(ctx, pr, body)
	}
	return nil
}

func (h *Handler) nagReviewerIfSlow(ctx context.Context, pr *github.PullRequest) error {
	logger := h.logger.WithField("policy", "nagReviewerIfSlow")
	since := time.Since(pr.GetCreatedAt())

	// If PR has just been filed, we do not comment
	if since < policyNagReviewerThreshold {
		logger.Infoln("PR is too new, bailing")
		return nil
	}

	reviewers, err := h.getReviewers(ctx, pr)
	if err != nil {
		return errors.Wrap(err, "issue getting reviewers")
	}

	reviews, err := h.getReviews(ctx, pr)
	if err != nil {
		return errors.Wrap(err, "issue getting reviews")
	}
	reviewerString := ""
	for _, reviewer := range reviewers.Users {
		reviewerString += "@"
		reviewerString += *reviewer.Login
		reviewerString += ", "
	}

	// Since reviews can be submitted by those that are not maintainers
	lastTimeReviewWasDoneByMaintainer := time.Time{}
	// var last_reviewer string
	for _, review := range reviews {
		// last_reviewer = *review.User.Login
		// TODO: prolly should explicitly check if the review was by a maintainer..?
		lastTimeReviewWasDoneByMaintainer = *review.SubmittedAt
	}

	comments, err := h.getComments(ctx, pr) // Go through comments and store the following
	if err != nil {
		return errors.Wrap(err, "issue getting PR comments")
	}
	lastTimeBotCommented := time.Time{}
	lastTimeSubmitterCommented := time.Time{}
	// var last_comment_author string = "No comments"

	// There should be a better way to do this. comments welcome
	for _, comment := range comments {
		commentUser := *comment.User.Login
		if commentUser == botUsername {
			lastTimeBotCommented = *comment.CreatedAt
		}
		if commentUser == *pr.User.Login {
			lastTimeSubmitterCommented = *comment.CreatedAt
		}
		// last_comment_author = comment_user
	}

	if time.Since(lastTimeReviewWasDoneByMaintainer) < policyNagReviewerThreshold {
		logger.Debugln(pr.GetNumber(), pr.GetTitle(), "Reviewer has recently submitted a review. Don't post again.")
		return nil
	}

	if time.Since(lastTimeSubmitterCommented) > policyNagSubmitterReviewCommentsThreshold && time.Since(lastTimeSubmitterCommented) > time.Since(lastTimeReviewWasDoneByMaintainer) {
		logger.Debugln(pr.GetNumber(), pr.GetTitle(), "Let's ping the submitter since the ball is in their court and a review has been done.")
		body := fmt.Sprintf("A review was posted by a maintainer. @%v, can you make the requested changes when you get a chance?", *pr.User.Login)
		return h.postComment(ctx, pr, body)
	}

	if time.Since(lastTimeBotCommented) < policyNagReviewerThreshold {
		logger.Debugln(pr.GetNumber(), pr.GetTitle(), "The bot has recently posted. Don't post again.")
		return nil
	}

	logger.Debugln(pr.GetNumber(), pr.GetTitle(), "If we got here, then we can remind the reviewer.")
	body := fmt.Sprintf("%vcan you review this PR when you get a chance?", reviewerString)
	return h.postComment(ctx, pr, body)
}

// Post a comment on the PR.
func (h *Handler) postComment(ctx context.Context, pr *github.PullRequest, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := h.client.GetIssuesService().CreateComment(ctx, githubOwner, githubRepo, pr.GetNumber(), comment)
	if err != nil {
		return errors.Wrap(err, "issue posting comment")
	}
	return nil
}

// Get reviews on the PR.
func (h *Handler) getReviews(ctx context.Context, pr *github.PullRequest) ([]*github.PullRequestReview, error) {
	opt := &github.ListOptions{}
	reviews, _, err := h.client.GetPullRequestsService().ListReviews(ctx, githubOwner, githubRepo, pr.GetNumber(), opt)
	if err != nil {
		return reviews, errors.Wrap(err, "issue getting reviews")
	}
	return reviews, nil
}

// Get comments on the PR. From GitHub API docs:
// "Comments on pull requests can be managed via the Issue Comments API."
func (h *Handler) getComments(ctx context.Context, pr *github.PullRequest) ([]*github.IssueComment, error) {
	comments, _, err := h.client.GetIssuesService().ListComments(ctx, githubOwner, githubRepo, pr.GetNumber(), nil)
	if err != nil {
		return comments, errors.Wrap(err, "issue getting PR comments")
	}
	return comments, nil
}

// Get requested reviewers on the PR.
func (h *Handler) getReviewers(ctx context.Context, pr *github.PullRequest) (*github.Reviewers, error) {
	opt := &github.ListOptions{}
	reviewers, _, err := h.client.GetPullRequestsService().ListReviewers(ctx, githubOwner, githubRepo, pr.GetNumber(), opt)
	if err != nil {
		return reviewers, errors.Wrap(err, "issue getting PR reviewers")
	}
	return reviewers, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
