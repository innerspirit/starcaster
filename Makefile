all:
	go build -ldflags="-s -w" . && ./starcaster