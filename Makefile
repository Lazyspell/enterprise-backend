SHELL_PATH = /bin/bash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/bash,/bin/bash)


run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

# =========================================================================
# Modules support

tidy: 
	go mod tidy
	go mod vendor
