build:
	bo build -o bin/server server.go

mac-build:
	GOOS=darwin GOARCH=amd64 go build -o bin/server server.go

windows-build:
	GOOS=windows GOARCH=amd64 go build -o bin/server.exe server.go

linux-build:
	GOOS=linux GOARCH=amd64 go build -o bin/server server.go

goose-sqlite-up:
	goose -dir migrations/sqlite sqlite3 fiber-petclinic-db.sqlite up

goose-sqlite-down:
	goose -dir migrations/sqlite sqlite3 fiber-petclinic-db.sqlite down

run:
	go run server.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy