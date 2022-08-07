APP_NAME := raphanus

test:
	go test -short -cover -v ./...

test-all:
	go test -cover -v ./...

test-race:
	go test -race -v ./...

lint:
	go vet ./...
	golangci-lint run

run-server:
	go run server/*.go

run-client-example:
	go run client/examples/simple.go

run-benchmark:
	go test -short -benchtime 5s -benchmem -bench .

watch-and-restart-server:
	reflex -s make run-server

build-docker-image:
	docker run --rm -v $$PWD:/go/src/$(APP_NAME) -w /go/src/$(APP_NAME) golang:alpine sh -c "apk add --no-cache git && go get ./... && go build -ldflags='-w -s' -o $(APP_NAME)-server ./server"
	docker build -t msoap/$(APP_NAME):latest .
	rm -f $(APP_NAME)-server
