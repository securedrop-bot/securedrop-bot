default:
	GOOS=linux GOARCH=386 go build
	docker build -t gcr.io/jen-personal/securedrop-bot .
	rm -f securedrop-bot

push:
	gcloud docker -- push gcr.io/jen-personal/securedrop-bot
