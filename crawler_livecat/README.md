# Crawler

- [Prerequisite](#prerequisite)
- [Build](#build)
- [Run](#run)
- [GoTest](#gotest)
  

## Prerequisite

### Install

- [Go 1.12 version](https://golang.org/dl/) - build middleware

### Golang

Use Go 1.12 or higher, add GOCACHE to `~/.bashrc`

```sh
sudo vim ~/.bashrc
export GOCACHE="/path/you/like/.cache/go-build"
```

Use Go 1.11 or below, delete `go.sum`, otherwise it will result error in verifying checksum.

```sh
rm go.sum
go get ./...
```

## Build

Install dependency

```sh
go mod download
```

Build

```sh
go build main.go
```

## Run

Run on talnet 120.126.16.88

```sh
cd /home/user/crawler/livecat/crawler_livecat/program/youtube
sudo ./youtube
cd /home/user/crawler/livecat/crawler_livecat/program/watermelon
sudo ./watermelon
cd /home/user/crawler/livecat/crawler_livecat/program/seventeen
sudo ./seventeen
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
go generate ./...
```

### Test

```sh
go test ./...
```