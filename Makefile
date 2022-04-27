ifdef OS
	RM = del /Q
	FixPath = $(subst /,\,$1)
	CLS = cls
else
   ifeq ($(shell uname), Linux)
		RM = rm -f
		FixPath = $1
		CLS = clear
   endif
endif

all:
	@echo make usage:
	@echo make run
	@echo make build
	@echo make test
	@echo make docker_run

run: build
	$(CLS) && $(call FixPath, web/backend)

build:
	cd web && go build .

test:
	$(CLS) && cd web && go test -v

tidy:
	$(CLS) && cd web && go mod tidy

docker_run:
	$(CLS) && cd web && docker-compose up --build

.PHONY: clean
clean:
	$(RM) $(call FixPath, web/backend*)
	$(RM) $(call FixPath, web/*.db)
	$(RM) $(call FixPath, *.db)
