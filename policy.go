package securedropbot

// Policy represents a bot policy.
type Policy struct {
	// InactivityThreshold is the waiting period between PR activity and bot action.
	InactivityThreshold time.Duration

	// AlreadyCommentedPattern is a regexp that determines if the bot has already carried out a policy (within AlreadyCommentedThreshold).
	AlreadyCommentedPattern *regexp.Regexp

	// AlreadyCommentedThreshold is the time span in which a previous comment is considered current.
	AlreadyCommentedThreshold time.Duration

	// Do carries out the policy if it is appropriate to do so.
	Do func(ctx context.Context, pr *github.PullRequest) error
}
