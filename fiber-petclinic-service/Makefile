build:
	go build -o bin/fiber-petclinic-service server.go

mock:
	mockery --all

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin//fiber-petclinic-service main.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin//fiber-petclinic-service.exe main.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin//fiber-petclinic-service main.go

goose-sqlite-up:
	goose -dir migrations/sqlite sqlite3 pet-clinic.sqlite up

goose-sqlite-down:
	goose -dir migrations/sqlite sqlite3 pet-clinic.sqlite down

test-coverage:
	go test -coverprofile=coverage.out ./...

view-test-coverage:
	go tool cover -html=coverage.out -o coverage.html

run:
	go run cmd/main.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy