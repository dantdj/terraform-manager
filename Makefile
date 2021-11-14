# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/cmd:
.PHONY: run/app
run/app:
	go run 

# ==================================================================================== #
# BUILD
# ==================================================================================== #

git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.version=${git_description}'

.PHONY: build/app
build/app:
	@echo "Building app..."
	GOOS=darwin GOARCH=arm64 go build -ldflags=${linker_flags} -o="./bin/darwin_arm64/tfm"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
.PHONY: audit
audit:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	@echo "Running tests..."
	go test -race -vet=off ./...