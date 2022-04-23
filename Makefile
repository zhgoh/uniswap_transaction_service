# This makefile is only tested in Windows 10
all:
	@echo make usage:
	@echo make run
	@echo make build

run:
	cd web && go run main.go

build:
	cd web && go build main.go

test:
	cd web && go test
