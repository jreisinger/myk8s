build: test
	go build
test:
	go test -race ./...
build-linux: test
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"