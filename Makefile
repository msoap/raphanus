test:
	go test -v ./...

lint:
	go vet ./...
	golint ./...
	errcheck ./...

gometalinter:
	gometalinter --vendor --cyclo-over=20 --line-length=150 --dupl-threshold=150 --min-occurrences=2 --enable=misspell --deadline=10m ./...

run-server:
	go run server/*.go

run-client-example:
	go run client/examples/simple.go

watch-and-restart-server:
	reflex -s make server-run
