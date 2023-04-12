FROM golang:1.13-alpine

ADD . /src
WORKDIR /src

RUN go build -o /src/walletApp main.go

ENTRYPOINT ./walletApp
