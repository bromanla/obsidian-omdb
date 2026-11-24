AIR_VERSION := v1.67.4
AIR_DIR := .air/bin/$(AIR_VERSION)
AIR := $(AIR_DIR)/air

.PHONY: help dev fmt

help:
	@printf '%s\n' \
		'make dev               Run the bot with live reload' \
		'make fmt               Format Go source files'

dev: $(AIR)
	./$(AIR) -c .air.toml

fmt:
	go fmt ./...

$(AIR):
	GOBIN="$(CURDIR)/$(AIR_DIR)" go install github.com/air-verse/air@$(AIR_VERSION)
