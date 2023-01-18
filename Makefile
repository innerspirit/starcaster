dev:
	go build -ldflags="-s -w" . && ./starcaster

win:
	go build -ldflags -H=windowsgui . && ./starcaster
