default:
	docker build --no-cache -t gcr.io/jen-personal/securedrop-bot .

push:
	gcloud docker -- push gcr.io/jen-personal/securedrop-bot
