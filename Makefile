pwd := $(shell pwd)
project := $(notdir $(pwd))


.PHONY:build linux win all clean upx install-lint lint callvis-install callvis proxy pprof


build:
	@go build -ldflags "-s -w"

test-all:


t1:
	@pwd | ./$(project)
t2:
	@echo "hello" | ./$(project)
t3:
	./$(project) -e ttc.mp
# flag can not be add at the last
# we could just create by ourself

linux:
	set  GOOS=linux
	go build -ldflags "-s -w" -o $(project)

win:
	@set GOOS=win
	@go build -ldflags "-s -w" -o $(project).exe
	@upx -9 $(project).exe

all:linux win

withoutwindow:
	@go build -ldflags "-s -w -H=windowsgui"


clean:
	-rm -f *.log
	-rm -f $(project)
	-rm -f *.exe
	-rm -f *.pprof
	-rm -f *.txt
	-rm -f *.mp

# tools-chain for golang


upx:build
	upx -9 $(project).exe

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	golangci-lint run

callvis-install:
	go install github.com/ofabry/go-callvis@master

callvis:
	go-callvis .

proxy:
	go env -w  GOPROXY=https://goproxy.io,direct

pprof:run
	go tool pprof -http=:8080 *.pprof