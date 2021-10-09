# build image
FROM golang:alpine as go_builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -v -trimpath -ldflags="-w -s" -o /go/bin/raphanus-server ./server/

# final image
FROM alpine

COPY --from=go_builder /go/bin/raphanus-server /app/raphanus-server
ENTRYPOINT ["/app/raphanus-server"]
CMD ["-address", ":8771"]
EXPOSE 8771
