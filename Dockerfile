FROM golang

ENV GOBIN /go/bin

RUN mkdir /securedrop-bot
RUN mkdir /go/src/securedrop-bot
ADD . /go/src/securedrop-bot
WORKDIR /go/src/securedrop-bot

RUN go get -u github.com/govend/govend
RUN govend -v

RUN go install

EXPOSE 8001
ENTRYPOINT ["/go/src/securedrop-bot/cmd/securedrop-bot/securedrop-bot"]
