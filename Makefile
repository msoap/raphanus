test:
	go test -short -cover -v ./...

test-all:
	go test -cover -v ./...

test-race:
	go test -race -v ./...

lint:
	go vet ./...
	golint ./...
	errcheck ./...

gometalinter:
	gometalinter --vendor --cyclo-over=20 --line-length=150 --dupl-threshold=150 --min-occurrences=2 --enable=misspell --deadline=10m --exclude=SA1022 ./...

run-server:
	go run server/*.go

run-client-example:
	go run client/examples/simple.go

run-benchmark:
	go test -short -benchtime 5s -benchmem -bench .

watch-and-restart-server:
	reflex -s make run-server

docker-build-image:
	rocker build --no-cache
	rm server/raphanus-server
