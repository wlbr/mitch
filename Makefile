LINKERFLAGS = -X main.Version=`git describe --tags` -X main.Githash=`git describe --always --long --dirty` -X main.Buildstamp=UTC`date -u '+%Y-%m-%d_%I:%M:%S%p'`
HOST = wlbr

all: build

build:
	go build -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go

run:
	go run -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go


deploy:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go
	rsync -a --progress  mitch $(HOST):./bin/
	ssh $(HOST) killall mitch
	ssh $(HOST) nohup bin/mitch &


dep:
	go get -u github.com/spf13/viper
	go get -u github.com/sbstjn/hanu