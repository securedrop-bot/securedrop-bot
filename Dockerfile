FROM golang

RUN go get -u github.com/govend/govend

COPY . /go/src/github.com/securedrop-bot/securedrop-bot
WORKDIR /go/src/github.com/securedrop-bot/securedrop-bot

RUN govend -v

RUN go install github.com/securedrop-bot/securedrop-bot/cmd/securedrop-bot

EXPOSE 8001

CMD ["securedrop-bot"]
