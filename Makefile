LINKERFLAGS = -X main.Version=`git describe --tags --always --long --dirty` -X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S_UTC'`

HOST = wlbr

all: clean build

.PHONY: clean
clean:
	rm -f mitch
	rm -rf bin/

build:
	go build -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go

run:
	go run -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go

debug:
	dlv debug cmd/mitch/mitch.go


rbuild: rbuild
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go

deploy: rbuild
	rsync -a --progress  mitch $(HOST):./bin/
	ssh $(HOST) killall mitch
	ssh $(HOST) nohup bin/mitch &

rstart:
	ssh $(HOST) nohup bin/mitch &

dep:
	go get -u github.com/spf13/viper
	go get -u github.com/nlopes/slack
