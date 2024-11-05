test:
	go test -race ./...
install: test
	go install
build: test
	go build
build-linux: test
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
