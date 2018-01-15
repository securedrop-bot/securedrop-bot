package securedropbot

import (
	"context"

	"github.com/google/go-github/github"
)

// pullRequestsService is an interface generated for "github.com/google/go-github/github".PullRequestsService.
type pullRequestsService interface {
	List(context.Context, string, string, *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
	ListReviewers(context.Context, string, string, int, *github.ListOptions) (*github.Reviewers, *github.Response, error)
	ListReviews(context.Context, string, string, int, *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error)
}
