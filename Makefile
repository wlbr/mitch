LINKERFLAGS = -X main.Version=`git describe --tags` -X main.Githash=`git describe --always --long --dirty` -X main.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`


all: build

build:
	go build -ldflags "$(LINKERFLAGS)" cmd/mitch/mitch.go

install:
	go install ./...

deploy:
	GOOS=linux GOARCH=amd64 go build cmd/mitch/mitch.go
	rsync -a --progress  mitch wlbr:./bin/
	ssh wlbr killall mitch
	ssh wlbr nohup bin/mitch &
