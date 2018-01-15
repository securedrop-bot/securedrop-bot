package securedropbot

import (
	"context"

	"github.com/google/go-github/github"
)

// repositoriesService is an interface generated for "github.com/google/go-github/github".RepositoriesService.
type repositoriesService interface {
	GetCombinedStatus(context.Context, string, string, string, *github.ListOptions) (*github.CombinedStatus, *github.Response, error)
}
