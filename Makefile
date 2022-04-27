# This makefile is only tested in Windows 10
all:
	@echo make usage:
	@echo make run
	@echo make build
	@echo make test
	@echo make docker_run

run: build
	clear
	.\web\backend.exe

build:
	cd web && go build .

test:
	clear && cd web && go test

tidy:
	clear && cd web && go mod tidy

docker_run:
	clear && cd web && docker-compose up --build
