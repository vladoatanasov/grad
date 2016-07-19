all: darwin linux64 linux32
darwin:
	go get && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o bin/grad_darwin
linux64:
	go get && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o bin/grad_linux_amd64
linux32:
	go get && GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-s -w -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=`git rev-parse HEAD`" -o bin/grad_linux86
clean:
	rm -fr bin/
buildclean: clean buildit
cleanbuild: clean buildit
test:
	go get && go test ./...
