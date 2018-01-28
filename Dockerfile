FROM alpine

RUN apk update
ADD cmd/securedrop-bot/securedrop-bot /securedrop-bot
ENTRYPOINT ["/securedrop-bot"]
