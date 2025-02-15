# Pet Clinic Restful API

## Goals
Build a Pet Clinic restful API using the libraries mentioned below.

## Libraries
* Dependency management: [go mod](https://blog.golang.org/using-go-modules)
* Application Configuration: [viper](https://github.com/spf13/viper)
* Routing framework: [gin gonic](https://github.com/gin-gonic/gin)
* Swagger:
  * Gin Swagger: [gin-swagger](https://github.com/swaggo/gin-swagger)
  * Go Swagger: [go-swagger](https://github.com/go-swagger/go-swagger)
* Actuator: [actuator]()
* Okta JWT: [okta](https://github.com/okta/okta-jwt-verifier-golang)
* Prometheus Client: [prometheus-client](https://github.com/prometheus/client_golang/prometheus/promhttp)
* Database: 
  * [postgresql](gorm.io/driver/postgres)
  * [sqlite3](gorm.io/driver/sqlite)
* Database ORM: [GORM](https://gorm.io/)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* Testing & Mock: [testify](https://github.com/stretchr/testify)
* Mockery: [mockery](https://github.com/vektra/mockery)

## Run HTTPS Server
To run the server with HTTPS, you need to generate a self-signed certificate and key. You can use the following command to generate the certificate and key.

```shell
openssl genrsa -des3 -out myCA.key 2048
openssl req -x509 -new -nodes -key myCA.key -sha256 -days 825 -out myCA.pem

-- convert passphrase to a key without passphrase, because SSL/TLS does not support passphrase
openssl rsa -in myCA.key -out myCA-without-passphrase.key
```

## App Health & Info Endpoints:

| Path     | Method | Description                   |
| :------- | :----- |:------------------------------|
| /health  | GET    | Check the service's heartbeat |
| /info    | GET    | Show the app info             |
| /metrics | GET    | Out of the box metrics        |

## API Endpoints:

| Path          | Method | Description           |
| :------------ | :----- | :-------------------- |
| /v1/vets/:id  | GET    | Get vetenarian by id  |
| /v1/pets/:id  | GET    | Get pet by id         |
| /v1/owner/:id | GET    | Get pet's owner by id |

---

## Mockery
Mockery is a tool used to generate mock implementations of Go interfaces. It is useful when you are writing unit tests for your code and you need to mock the dependencies of the code under test.

### Installing mockery
Install the mockery tool into your system using Homebrew and go get command.
```shell
brew install vektra/tap/mockery
````

Add vektra mockery to go mod 
```shell
go get github.com/vektra/mockery/v3
```

**Creating mock function/interface**<br/>
```yaml
with-expecter: true
mockname: "Mock{{.InterfaceName}}"
inpackage: true
case: underscore
```
Run this command to generate the mock files
```shell
mockery --all
```

## [Uber Gomock](https://github.com/uber-go/mock)
GoMock is a mocking framework for the Go programming language. It integrates well with Go's built-in testing package, but can be used in other contexts too.

### Installation
```bash
go install go.uber.org/mock/mockgen@latest
mockgen -version
```

## Steps to run the application

1. create a database and tables
```shell
goose -dir migrations/sqlite up
```

2. Run the application
```shell
go run cmd/server.go
```

## Swagger
Swagger is a tool that can help you design, build, document, and consume RESTful web services. It is a specification for describing, producing, consuming, and visualizing RESTful web services.

### Install Swagger
```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger Docs
Run this command in the root directory of the project.  It will search main.go file for the gin engine and generate the swagger docs.
```shell
swwag init
```
if you have a different main file and in subdirectory named cmd, you can specify it like this:

```shell
swag init -g cmd/server.go
```
After using swag init to generate Swagger 2.0 docs, import the following package in your main.go file:

```go
package main

import (
  _ "your_project/docs"
  swaggerFiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"
)
````

**_NOTE:_** You need to import `_ "your_project/docs"` to ensure that the docs are generated and available at runtime.







## Okta Authentication & Authorization

The Client Credentials flow is recommended for use in machine-to-machine authentication. Your application will need to securely store its Client ID and Secret and pass those to Okta in exchange for an access token. At a high-level, the flow only has two steps:

* Your application passes its client credentials to your Okta authorization server.
* If the credentials are accurate, Okta responds with an access token.

## Project Layout
The template project layout:

```
.
├── cmd                   main application
├── config                configuration files for different environments
│    └── api              app configuration                 
├── docs                  swagger docs files
├── errors                error types and handling
├── internal              private application & library
│    └── api              application API
│         ├── metrics     metrics endpoint
│         ├── health      health check endpoint
│         └── info        info endpoint
│         └── owner       owner endpoint
├── middleware            middleware
├── migrations            database migrations & schema
├── pkg                   shared libraries
│     ├── dbase           database server library
│     ├── da              http server library
│     └── infra           infra resouces
│          └── repository repository library      
└── testdata              test data scripts
```

## Application Architecture
Explain the design patterns, the libraries, and the frameworks used in the application.

### Entry Point
The server.go is the application entrypoint. It comprises two init functions. The first init function responsible
for loading the application configuration, logger, and gin engine. The second init function is responsible for
setting up the application routes and middleware.

### Configuration
The app.LoadConfig function is responsible for loading the application configuration in YAML format. It reads and
converts the yaml file into a struct using the viper library.  For example, it read this yaml file:

```yaml
appInfo:
  name: "Pet Clinic"
  version: "1.0.0"
  description: "Pet Clinic Restful"

database:
  postgres:
    driver: postgres
    dsn: "user=postgres password=mysecretpassword host=localhost port=5432 dbname=petclinic sslmode=disable"
  maxIdleConns: 0
  maxOpenConns: 5
  maxIdleTime: 60
```

and converts it into this struct:

**By my convention, I add postfix Config to the struct name to indicate its type as a configuration struct.**
```go
type AppInfoConfig struct {
    Name        string `yaml:"name"`
    Version     string `yaml:"version"`
    Description string `yaml:"description"`
}

// I don't specify the yaml because the variable names match to the yaml key name
type DatabaseConfig struct {
    Postgres     PostgresConfig
    MaxIdleConns int
    MaxOpenConns int
    MaxIdleTime  int
}

type PostgresConfig struct {
    Driver string
    Dsn    string
}
```


### Components Instantiation


### Inverting Dependencies



### Middleware



### Router Register



### Unit Tests



### Swagger


### A

### A

### A

### A
