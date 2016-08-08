test:
	go test -v

lint:
	go vet ./...
	golint ./...
	errcheck ./...

server-run:
	go run server/*.go

watch-and-restart-server:
	reflex -s make server-run
