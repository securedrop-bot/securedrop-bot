package securedropbot

import (
	"context"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

var failureStatus = "failure"
var botLogin = "securedrop-bot"
var fiveMinutesAgo = time.Now().Add(-5 * time.Minute)

var (
	botUser = &github.User{
		Login: &botLogin,
	}
	failureCombinedStatus = &github.CombinedStatus{
		State: &failureStatus,
		Statuses: []github.RepoStatus{
			github.RepoStatus{
				State: &failureStatus,
			},
		},
	}
	recentBotComments = []*github.IssueComment{
		&github.IssueComment{
			User:      botUser,
			CreatedAt: &fiveMinutesAgo,
		},
	}
)

// mocks
var (
	failingStatusRepository = &repositoriesServiceMock{
		GetCombinedStatusFunc: func(in1 context.Context, in2 string, in3 string, in4 string, in5 *github.ListOptions) (*github.CombinedStatus, *github.Response, error) {
			return failureCombinedStatus, nil, nil
		},
	}
	noRecentComments = &issuesServiceMock{
		CreateCommentFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return nil, nil, nil
		},
		ListCommentsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
			return []*github.IssueComment{}, nil, nil
		},
	}
	botRecentlyCommentedIssues = &issuesServiceMock{
		CreateCommentFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return nil, nil, nil
		},
		ListCommentsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
			return recentBotComments, nil, nil
		},
	}
)

func TestHandler_nagSubmitterIfFailed(t *testing.T) {
	logger := logrus.New()
	if testing.Verbose() {
		logger.SetLevel(logrus.DebugLevel)
	}

	ctx := context.Background()
	type args struct {
		ctx context.Context
		pr  *github.PullRequest
	}
	tests := []struct {
		name      string
		client    *fakeGithubClient
		args      args
		nComments int // how many comments we expect to be posted
		wantErr   bool
	}{
		{"failing-status", &fakeGithubClient{
			repositories: failingStatusRepository,
			pullRequests: &pullRequestsServiceMock{},
			issues:       noRecentComments,
		}, args{ctx: ctx}, 1, false},
		{"failing-status-with-comment-present", &fakeGithubClient{
			repositories: failingStatusRepository,
			pullRequests: &pullRequestsServiceMock{},
			issues:       botRecentlyCommentedIssues,
		}, args{ctx: ctx}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				logger: logger,
				client: tt.client,
			}
			if err := h.nagSubmitterIfFailed(tt.args.ctx, tt.args.pr); (err != nil) != tt.wantErr {
				t.Errorf("Handler.nagSubmitterIfFailed() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := len(tt.client.issues.calls.CreateComment)
			if got != tt.nComments {
				t.Errorf("Handler.nagSubmitterIfFailed, got %d comments; want %d", got, tt.nComments)
			}
		})
	}
}
