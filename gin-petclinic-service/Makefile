build:
	bo build -o bin/server server.go

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin/server server.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe server.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin/server server.go

mock:
	mockery

swag:
	swag i -g cmd/server.go

goose-sqlite-up:
	goose -dir migrations/sqlite sqlite3 pet-clinic.sqlite up

goose-sqlite-down:
	goose -dir migrations/sqlite sqlite3 pet-clinic.sqlite down

goose-postgres-up:
	goose -dir migrations/postgres postgres  up

goose-postgres-down:
	goose -dir migrations/postgres postgres pet-clinic.sqlite down

postgres-start:
	docker run --name local-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:latest

postgres-stop:
	docker stop local-postgres || true

test-coverage:
	go test -coverprofile=coverage.out ./...

view-test-coverage:
	go tool cover -html=coverage.out -o coverage.html

live-reload:
	air -d -c ./air/.air.linux.conf

run:
	go run cmd/server.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy
	
docker-build:
	docker build -t gin-petclinic-service:latest .