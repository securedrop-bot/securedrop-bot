package securedropbot

type fakeGithubClient struct {
	repositories *repositoriesServiceMock
	pullRequests *pullRequestsServiceMock
	issues       *issuesServiceMock
}

// GetRepositoriesService returns a repositories service.
func (c *fakeGithubClient) GetRepositoriesService() repositoriesService {
	return c.repositories
}

// GetPullRequestsService returns a pull requests service.
func (c *fakeGithubClient) GetPullRequestsService() pullRequestsService {
	return c.pullRequests
}

// GetIssuesService returns a issues service.
func (c *fakeGithubClient) GetIssuesService() issuesService {
	return c.issues
}
