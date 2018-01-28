package securedropbot

import (
	"context"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
)

var (
	failureStatus       = "failure"
	botLogin            = "securedrop-bot"
	fiveMinutesAgo      = time.Now().Add(-5 * time.Minute)
	threeDaysAgo        = time.Now().Add(-72 * time.Hour)
	approvedStatus      = "APPROVED"
	requestChangeStatus = "CHANGES_REQUESTED"
	genericPRNumber     = 1
	userLogin           = "redshiftzero"
	prSubmitterLogin    = "contributorname"
)

var (
	botUser = &github.User{
		Login: &botLogin,
	}
	reviewerUser = &github.User{
		Login: &userLogin,
	}
	prSubmitterUser = &github.User{
		Login: &prSubmitterLogin,
	}
	genericReviewers      = &github.Reviewers{}
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
	recentReviewerComments = []*github.IssueComment{
		&github.IssueComment{
			User:      reviewerUser,
			CreatedAt: &fiveMinutesAgo,
		},
	}
	approvedReviews = []*github.PullRequestReview{
		&github.PullRequestReview{
			State:       &approvedStatus,
			User:        reviewerUser,
			SubmittedAt: &threeDaysAgo,
		},
	}
	changesRequestedReviews = []*github.PullRequestReview{
		&github.PullRequestReview{
			State:       &requestChangeStatus,
			User:        reviewerUser,
			SubmittedAt: &threeDaysAgo,
		},
	}
	recentReview = []*github.PullRequestReview{
		&github.PullRequestReview{
			State:       &requestChangeStatus,
			User:        reviewerUser,
			SubmittedAt: &fiveMinutesAgo,
		},
	}
	genericPullRequest = &github.PullRequest{
		Number: &genericPRNumber,
		User:   prSubmitterUser,
	}
	recentPullRequest = &github.PullRequest{
		Number:    &genericPRNumber,
		CreatedAt: &fiveMinutesAgo,
		User:      prSubmitterUser,
	}
)

// mocks

func noRecentComments() *issuesServiceMock {
	return &issuesServiceMock{
		CreateCommentFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return nil, nil, nil
		},
		ListCommentsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
			return []*github.IssueComment{}, nil, nil
		},
	}
}

func failingStatusRepository() *repositoriesServiceMock {
	return &repositoriesServiceMock{
		GetCombinedStatusFunc: func(in1 context.Context, in2 string, in3 string, in4 string, in5 *github.ListOptions) (*github.CombinedStatus, *github.Response, error) {
			return failureCombinedStatus, nil, nil
		},
	}
}
func botRecentlyCommentedIssues() *issuesServiceMock {
	return &issuesServiceMock{
		CreateCommentFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return nil, nil, nil
		},
		ListCommentsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
			return recentBotComments, nil, nil
		},
	}
}
func reviewerRecentlyCommentedIssues() *issuesServiceMock {
	return &issuesServiceMock{
		CreateCommentFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return nil, nil, nil
		},
		ListCommentsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
			return recentReviewerComments, nil, nil
		},
	}
}
func prIsApproved() *pullRequestsServiceMock {
	return &pullRequestsServiceMock{
		ListReviewsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
			return approvedReviews, nil, nil
		},
	}
}
func prIsNotApproved() *pullRequestsServiceMock {
	return &pullRequestsServiceMock{
		ListReviewersFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			return genericReviewers, nil, nil
		},
		ListReviewsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
			return changesRequestedReviews, nil, nil
		},
	}
}
func prIsRecentlyReviewed() *pullRequestsServiceMock {
	return &pullRequestsServiceMock{
		ListReviewersFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			return genericReviewers, nil, nil
		},
		ListReviewsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
			return recentReview, nil, nil
		},
	}
}
func prHasNoReviews() *pullRequestsServiceMock {
	return &pullRequestsServiceMock{
		ListReviewersFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) (*github.Reviewers, *github.Response, error) {
			return genericReviewers, nil, nil
		},
		ListReviewsFunc: func(in1 context.Context, in2 string, in3 string, in4 int, in5 *github.ListOptions) ([]*github.PullRequestReview, *github.Response, error) {
			return nil, nil, nil
		},
	}
}

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
			repositories: failingStatusRepository(),
			pullRequests: &pullRequestsServiceMock{},
			issues:       noRecentComments(),
		}, args{ctx: ctx}, 1, false},
		{"failing-status-with-comment-present", &fakeGithubClient{
			repositories: failingStatusRepository(),
			pullRequests: &pullRequestsServiceMock{},
			issues:       botRecentlyCommentedIssues(),
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

func TestHandler_nagMaintainerForMerge(t *testing.T) {
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
		{"pr-is-not-approved", &fakeGithubClient{
			pullRequests: prIsNotApproved(),
			issues:       noRecentComments(),
		}, args{ctx: ctx, pr: genericPullRequest}, 0, false},
		{"pr-is-approved-no-comments", &fakeGithubClient{
			pullRequests: prIsApproved(),
			issues:       noRecentComments(),
		}, args{ctx: ctx, pr: genericPullRequest}, 1, false},
		{"pr-is-approved-bot-already-commented", &fakeGithubClient{
			pullRequests: prIsApproved(),
			issues:       botRecentlyCommentedIssues(),
		}, args{ctx: ctx, pr: genericPullRequest}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				logger: logger,
				client: tt.client,
			}
			if err := h.nagMaintainerForMerge(tt.args.ctx, tt.args.pr); (err != nil) != tt.wantErr {
				t.Errorf("Handler.nagMaintainerForMerge() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := len(tt.client.issues.calls.CreateComment)
			if got != tt.nComments {
				t.Errorf("Handler.nagMaintainerForMerge, got %d comments; want %d", got, tt.nComments)
			}
		})
	}
}

func TestHandler_nagReviewerIfSlow(t *testing.T) {
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
		{"pr-is-too-new", &fakeGithubClient{
			pullRequests: prIsNotApproved(),
			issues:       noRecentComments(),
		}, args{ctx: ctx, pr: recentPullRequest}, 0, false},
		{"review-recently-submitted", &fakeGithubClient{
			pullRequests: prIsRecentlyReviewed(),
			issues:       noRecentComments(),
		}, args{ctx: ctx, pr: genericPullRequest}, 0, false},
		{"nag-submitter-to-make-requested-changes", &fakeGithubClient{
			pullRequests: prIsNotApproved(),
			issues:       reviewerRecentlyCommentedIssues(),
		}, args{ctx: ctx, pr: genericPullRequest}, 1, false},
		{"if-bot-recently-commented-do-not-nag", &fakeGithubClient{
			pullRequests: prIsNotApproved(),
			issues:       botRecentlyCommentedIssues(),
		}, args{ctx: ctx, pr: genericPullRequest}, 0, false},
		{"bot-should-nag-reviewers", &fakeGithubClient{
			pullRequests: prHasNoReviews(),
			issues:       noRecentComments(),
		}, args{ctx: ctx, pr: genericPullRequest}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				logger: logger,
				client: tt.client,
			}
			if err := h.nagReviewerIfSlow(tt.args.ctx, tt.args.pr); (err != nil) != tt.wantErr {
				t.Errorf("Handler.nagReviewerIfSlow() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := len(tt.client.issues.calls.CreateComment)
			if got != tt.nComments {
				t.Errorf("Handler.nagReviewerIfSlow, got %d comments; want %d", got, tt.nComments)
			}
		})
	}
}
