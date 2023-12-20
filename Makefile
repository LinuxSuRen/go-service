test:
	go fmt ./...
	go test ./...
	GOOS=windows go build .
	GOOS=darwin go build .
	GOOS=linux go build .
build:
	go build -o bin/go-svc demo/main.go
install: build
	sudo install bin/go-svc /usr/local/bin
hd:
	curl https://linuxsuren.github.io/tools/install.sh|bash
init-env: hd
	hd i cli/cli
	gh extension install linuxsuren/gh-dev
