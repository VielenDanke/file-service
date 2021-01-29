FROM golang:alpine as build-env

ENV DB_URL=postgres://user:userpassword@172.17.0.2:5432/file_service_db?sslmode=disable
ENV SERVER_PORT=:4545
ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /file_service

WORKDIR /file_service

COPY . .

RUN go build -o file_service .

CMD ./file_service