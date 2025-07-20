build:
	go build -o bin/fiber-petclinic-service server.go

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin//fiber-petclinic-service main.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin//fiber-petclinic-service.exe main.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin//fiber-petclinic-service main.go

run:
	go run cmd/main.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy