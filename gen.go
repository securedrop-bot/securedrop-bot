//go:generate moq -out mock_issues_service_test.go . issuesService
//go:generate moq -out mock_repositories_service_test.go . repositoriesService
//go:generate moq -out mock_pull_requests_service_test.go . pullRequestsService
package securedropbot
