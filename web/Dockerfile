# syntax=docker/dockerfile:1
# https://docs.docker.com/language/golang/build-images/

FROM golang:1.18.1-alpine3.15

ARG buildtime_variable=default_value  # <- this one's new
ENV etherscan_api=$buildtime_variable # we reference it directly

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /backend

EXPOSE 5050

CMD ["/backend"]
