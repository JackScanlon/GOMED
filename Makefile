GOOS = linux
CGO_ENABLED = 0
DEFAULT_VER = 0.0.1
DEFAULT_CMD = build

dev: build-dev run-dev

run-dev:
	$(eval ARGS := $(DEFAULT_CMD))
	@./out/bin/snomed $(ARGS)

build-dev:
	$(eval VERSION = $(shell git describe --abbrev=0 2>/dev/null || echo ${DEFAULT_VER}))
	$(eval LDFLAGS = -ldflags "-w -s -X snomed/src/shared.Version=$(VERSION)")
	$([ -d out/bin ] || out/bin)
	@go build -a ${LDFLAGS} -o out/bin/snomed src/main.go

clean:
	@rm -r out/
