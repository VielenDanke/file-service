FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /file_service

WORKDIR /file_service

COPY . .

RUN go build -o file_service .

CMD ./file_service