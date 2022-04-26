# syntax=docker/dockerfile:1
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18.1-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /backend

EXPOSE 5050

CMD ["/backend"]
