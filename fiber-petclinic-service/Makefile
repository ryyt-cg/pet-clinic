build:
	go build -o bin/fiber-petclinic-service cmd/server.go

mock:
	mockery --all

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin/fiber-petclinic-service cmd/server.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin/fiber-petclinic-service.exe cmd/server.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin/fiber-petclinic-service cmd/server.go

swag:
	swag i -g cmd/server.go

goose-sqlite-up:
	goose -env .env.sqlite up

goose-sqlite-down:
	goose -dir migrations/sqlite sqlite3 pet-clinic.sqlite down

goose-postgres-up:
	goose -env .env.postgres up

goose-postgres-down:
	goose -dir migrations/postgres postgres pet-clinic.sqlite down

postgres-start:
	docker run --name local-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:latest

postgres-stop:
	docker stop local-postgres || true

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

view-test-coverage:
	go tool cover -html=coverage.out -o coverage.html

live-reload:
	air -c ./air/.air.linux.conf

run:
	go run cmd/main.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy

docker-build:
	docker build -t fiber-petclinic-service:latest .

docker-run:
	docker rm -f fiber-petclinic-service || true
	docker run -d -p 8091:8091 --name fiber-petclinic-service fiber-petclinic-service:latest