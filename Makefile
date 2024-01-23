PROJECT := sanic

.DEFAULT_GOAL := help

.PHONY: mpd run build tidy verify test cert container help

mpd:  ## Run mpd test instance
	mkdir -p /tmp/${PROJECT}/{music,playlists}
	touch /tmp/${PROJECT}/mpd_db
	mpd --no-daemon ./mpd.conf

run: build  ## Run project
	./${PROJECT}

build:  ## Compile project
	go build -v -o ${PROJECT}

update:  ## Update go dependencies
	go get -u

tidy:  ## Add missing and remove unused modules
	go mod tidy

verify:  ## Verify dependencies have expected content
	go mod verify

test:  ## Run tests
	go test ./...

cert:  ## Create https certificate for local testing
	go run $$GOROOT/src/crypto/tls/generate_cert.go --host localhost

build-container:  ## Build container image
	podman build --tag ${PROJECT}:latest .

run-container: build-container  ## Run container image
	podman run --rm --volume ./config.ini:/config.ini --publish-all ${PROJECT}:latest

help: ## Display this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

