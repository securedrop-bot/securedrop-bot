package securedropbot

import (
	"context"

	"github.com/google/go-github/github"
)

// issuesService is an interface generated for "github.com/google/go-github/github".IssuesService.
type issuesService interface {
	CreateComment(context.Context, string, string, int, *github.IssueComment) (*github.IssueComment, *github.Response, error)
	ListComments(context.Context, string, string, int, *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error)
}
