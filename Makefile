# This makefile is only tested in Windows 10
all:
	@echo make usage:
	@echo make run
	@echo make build

run: build
	cls
	.\web\main.exe

build:
	cd web && go build main.go etherscan.go

test:
	cls && cd web && go test
