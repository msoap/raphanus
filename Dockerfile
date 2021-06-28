# build image
FROM golang:alpine as go_builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0
RUN go install -v -ldflags="-w -s" github.com/msoap/raphanus/server

# final image
FROM alpine

COPY --from=go_builder /go/bin/server /app/raphanus-server
ENTRYPOINT ["/app/raphanus-server"]
CMD ["-address", ":8771"]
EXPOSE 8771
