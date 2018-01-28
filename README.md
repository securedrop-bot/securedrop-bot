# SecureDrop Bot

## Deployment

```
$ kubectl create secret generic securedrop-bot-github --from-literal=api_token='TOKEN_GOES_HERE'
$ kubectl create -f k8s.yaml
```
