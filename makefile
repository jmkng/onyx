.PHONY: build

BINARY_NAME=onyx

build:
	GOARCH=arm64 GOOS=darwin go build -o ./build/${BINARY_NAME}-darwin-arm *.go; \
	GOARCH=amd64 GOOS=darwin go build -o ./build/${BINARY_NAME}-darwin *go; \
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME}-linux *.go; \
	GOARCH=amd64 GOOS=windows go build -o ./build/${BINARY_NAME}-windows *.go; \

clean:
	go clean; \
	rm ./build/${BINARY_NAME}-darwin-arm; \
	rm ./build/${BINARY_NAME}-darwin; \
	rm ./build/${BINARY_NAME}-linux; \
	rm ./build/${BINARY_NAME}-windows; \