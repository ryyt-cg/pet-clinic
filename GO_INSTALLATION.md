# Golang Environment Setup

## Mac OS

### Go Version Manager

There are two popular go version managers:

1) [Simple go version manager](https://github.com/stefanmaric/g)
2) Go Version Manager [GVM](https://github.com/moovweb/gvm)

I personally chose the first one, simple version. I simply need a feature to install multiple versions of go and switch
between them.

## Windows

### Go Version Manager

Will update soon.

## Linux

### Go Version Manager

Will update soon.

## Install Go Version Manager

```shell
curl -sSL https://git.io/g-install | sh -s
```

### Usage

```
  Usage: g [COMMAND] [options] [args]

  Commands:

    g                         Open interactive UI with downloaded versions
    g install latest          Download and set the latest go release
    g install <version>       Download and set go <version>
    g download <version>      Download go <version>
    g set <version>           Switch to go <version>
    g run <version>           Run a given version of go
    g which <version>         Output bin path for <version>
    g remove <version ...>    Remove the given version(s)
    g prune                   Remove all versions except the current version
    g list                    Output downloaded go versions
    g list-all                Output all available, remote go versions
    g self-upgrade            Upgrades g to the latest version
    g help                    Display help information, same as g --help

  Options:

    -h, --help                Display help information and exit
    -v, --version             Output current version of g and exit
    -q, --quiet               Suppress almost all output
    -c, --no-color            Force disabled color output
    -y, --non-interactive     Prevent prompts
    -o, --os                  Override operating system
    -a, --arch                Override system architecture
    -u, --unstable            Include unstable versions in list
```



[fiber v2](fiber-petclinic-service/README.md) |
[fiber v3](fiber3-petclinic-service/README.md) |
[gin gonic](gin-petclinic-service/README.md) |