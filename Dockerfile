FROM golang:1.11-alpine

WORKDIR /go/src/app

RUN apk update && apk upgrade && \
  apk add --no-cache bash git

RUN go get github.com/golang/dep/cmd/dep

ADD . /go/src/app/

RUN dep ensure