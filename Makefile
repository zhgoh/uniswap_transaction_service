# This makefile is only tested in Windows 10
all:
	@echo make usage:
	@echo make run
	@echo make build

run: build
	cls
	.\web\backend.exe

build:
	cd web && go build .

test:
	cls && cd web && go test

tidy:
	cls && cd web && go mod tidy

docker_build:
	sudo docker build --build-arg buildtime_variable=a_value --tag backend .

docker_run:
	sudo docker run --rm -p 5000:5050 backend
