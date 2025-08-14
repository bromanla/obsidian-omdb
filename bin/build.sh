#!/bin/bash
# Build the binary

source "$(dirname "$0")/common.sh"

info "Building ${SERVICE_NAME}..."
go build -o ./bin/omdb-bot ./cmd
success "Build complete: ./bin/${SERVICE_NAME}"