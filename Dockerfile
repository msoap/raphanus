# build image
FROM golang:alpine as go_builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0
RUN go get -v github.com/msoap/raphanus/...
RUN cd /go/src/github.com/msoap/raphanus && go install -a -v -ldflags="-w -s" ./...

# final image
FROM alpine

COPY --from=go_builder /go/bin/server /app/raphanus-server
ENTRYPOINT ["/app/raphanus-server"]
CMD ["-address", ":8771"]
EXPOSE 8771
