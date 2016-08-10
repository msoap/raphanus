test:
	go test -v ./...

lint:
	go vet ./...
	golint ./...
	errcheck ./...

run-server:
	go run server/*.go

run-client-example:
	go run client/examples/simple.go

watch-and-restart-server:
	reflex -s make server-run
