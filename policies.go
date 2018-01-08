package securedropbot

var (
	nagCommitAuthorsOnFailure = Policy{
		InactivityThreshold:     2 * time.Hour,
		AlreadyCommentedPattern: regexp.MustCompile("looks like there was a test failure"),
		Do: func(ctx context.Context, logger logrus.FieldLogger, client *github.Client, pr *github.PullRequest) error {
			// TODO implement
			return nil
		},
	}
)
