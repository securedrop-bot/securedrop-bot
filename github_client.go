package securedropbot

import "github.com/google/go-github/github"

type githubClient interface {
	GetRepositoriesService() repositoriesService
	GetPullRequestsService() pullRequestsService
	GetIssuesService() issuesService
}

// realGithubClient satisfies the githubClient interface with a concrete github.Client
type realGithubClient struct {
	*github.Client
}

// GetRepositoriesService returns a repositories service.
func (c *realGithubClient) GetRepositoriesService() repositoriesService {
	return c.Repositories
}

// GetPullRequestsService returns a pull requests service.
func (c *realGithubClient) GetPullRequestsService() pullRequestsService {
	return c.PullRequests

}

// GetIssuesService returns a issues service.
func (c *realGithubClient) GetIssuesService() issuesService {
	return c.Issues
}
