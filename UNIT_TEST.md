# Test Driven Development (TDD)

## [Mockery](https://vektra.github.io/mockery/latest/)

Mockery is a tool used to generate mock implementations of Go interfaces. It is useful when you are writing unit tests
for your code and you need to mock the dependencies of the code under test.

### Installing mockery

Install the mockery tool into your system using Homebrew and go get command.

```shell
brew install mockery
brew upgrade mockery
````

Go install the version of mockery

```shell
go install github.com/vektra/mockery/v3@v3.5.1
```

Add vektra mockery to go mod

```bash
go get github.com/vektra/mockery/v3
```

**Initialize mockery**<br/>

```bash
mockery init <module-name>
```

This command will create .mockery.yaml file in the root of your module. This file contains the configuration for
mockery.

Run this command to generate the mock files

```shell
mockery
```

[fiber v2](fiber-petclinic-service/README.md) |
[fiber v3](fiber3-petclinic-service/README.md) |
[gin gonic](gin-petclinic-service/README.md) |

