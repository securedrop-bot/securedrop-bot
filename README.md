# SecureDrop Bot

[![CircleCI](https://circleci.com/gh/securedrop-bot/securedrop-bot.svg?style=svg)](https://circleci.com/gh/securedrop-bot/securedrop-bot)
[![codecov](https://codecov.io/gh/securedrop-bot/securedrop-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/securedrop-bot/securedrop-bot)

## Description

I remind contributors to:

* Make changes to a PR once maintainer(s) have reviewed it.
* Investigate test failures.  

I remind maintainers to:

* Review PRs if more review is needed (from the requested reviewers).
* Merge a PR once it is approved.

## Deployment

```
$ kubectl create secret generic securedrop-bot-github --from-literal=api_token='TOKEN_GOES_HERE'
$ kubectl create -f k8s.yaml
```
