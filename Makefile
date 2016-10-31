all: build

build:
	go build -ldflags "-X main.xyz=abc" cmd/mitch/mitch.go

install:
	go install ./...

deploy:
	GOOS=linux GOARCH=amd64 go build cmd/mitch/mitch.go
	rsync -a --progress  mitch wlbr:./bin/
	ssh wlbr killall mitch
	ssh wlbr nohup bin/mitch &
	#ssh wlbr MITCH_SLACK_TOKEN=xoxb-95885393046-rzAjVkPcv6TYycbrr4gM0PAT  nohup bin/mitch &
