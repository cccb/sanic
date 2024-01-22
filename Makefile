mpd:
	mkdir -p /tmp/sanic/{music,playlists}
	touch /tmp/sanic/mpd_db
	mpd --no-daemon ./mpd.conf

run: build
	./server

build:
	go build -v -o sanic

tidy:  ## add missing and remove unused modules
	go mod tidy

verify:  ## verify dependencies have expected content
	go mod verify

test:  ## run tests
	go test ./...

cert:  ## create https certificate for testing
	go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost

container:
	podman build -t dotfiles:latest .

