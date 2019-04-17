# Crawler

- [Prerequisite](#prerequisite)
- [Build](#build)
- [GoTest](#gotest)
  

## Prerequisite

### Install

- cmake、gcc、g++ - build core library
- [Go 1.12 version](https://golang.org/dl/) - build middleware

### Golang

Use Go 1.12 or higher, add GOCACHE to `~/.bashrc`

```sh
sudo vim ~/.bashrc
export GOCACHE="/path/you/like/.cache/go-build"
```

Use Go 1.11 or below, delete `go.sum`, otherwise it will result error in verifying checksum.

```sh
cd src
rm go.sum
go get ./...
```

## Build

Install dependency

```sh
cd src
go mod download
```


## GoTest

### Generate GoMock file

Install gomock and mockgen

```sh
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen
```

Generate

```sh
cd src/
go generate ./...
```

### Test

```sh
cd src/
go test ./...
```