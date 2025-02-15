build:
	bo build -o bin/server server.go

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin/server server.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe server.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin/server server.go

run:
	go run server.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy