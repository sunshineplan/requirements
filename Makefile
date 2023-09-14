.PHONY: build test run clean

build:
	npm install
	npm run build
	GOARCH=arm64 GOOS=darwin go build -ldflags "-s -w" -o requirements-darwin-arm64
	GOARCH=amd64 GOOS=darwin go build -ldflags "-s -w" -o requirements-darwin-amd64
	GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o requirements-linux-amd64
	GOARCH=amd64 GOOS=windows go build -ldflags "-s -w" -o requirements-windows-amd64.exe

test:
	npm install
    npm run check
    npm run build
	go build -v ./...
	go clean

run:
	npm install
	npm run build
	go build -ldflags "-s -w"
	./requirements

clean:
	go clean
	rm -r dist
	rm -r node_modules
