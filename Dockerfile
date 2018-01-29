FROM golang

ENV GOBIN /go/bin

ADD . /go/src/github.com/tmc/securedrop-bot
WORKDIR /go/src/github.com/tmc/securedrop-bot

RUN go get -u github.com/govend/govend
RUN govend -v

RUN go install cmd/securedrop-bot/main.go

EXPOSE 8001
CMD ["/go/bin/main"]
