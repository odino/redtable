FROM golang:1.22-alpine

RUN apk add --update bash

RUN go install github.com/cosmtrek/air@latest

WORKDIR /src
COPY go.mod /src
COPY go.sum /src
RUN go mod download -x

COPY . /src

ENV BIGTABLE_EMULATOR_HOST=bigtable:8086

CMD air