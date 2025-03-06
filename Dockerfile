FROM golang:1.24-alpine3.21

WORKDIR /muslib

COPY . ./

RUN go mod tidy 