# vim: set ft=make ts=8 noet

.PHONY: help
help:	### Get some help about make targets.
### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

APP_NAME=tcp_to_mqtt

.PHONY: lint
lint: ### Run linter.
	golangci-lint run

build: ### Build.
	go build -a -o bin/tcp_to_mqtt ./main.go

serve: build ### Run cli application
	./bin/tcp_to_mqtt serve