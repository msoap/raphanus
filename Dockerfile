# build image
FROM --platform=$BUILDPLATFORM golang:alpine as go_builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git

ENV CGO_ENABLED=0
# GOARM=6 affects only "arm" builds
ENV GOARM=6
# "amd64", "arm64" or "arm" (--platform=linux/amd64,linux/arm64,linux/arm/v6)
ENV GOARCH=$TARGETARCH
ENV GOOS=linux
COPY . /src
WORKDIR /src
RUN echo "Building for $GOOS/$GOARCH"
RUN go build -v -trimpath -ldflags="-w -s" -o /go/bin/raphanus-server ./server/

# final image
FROM alpine

COPY --from=go_builder /go/bin/raphanus-server /app/raphanus-server
ENTRYPOINT ["/app/raphanus-server"]
CMD ["-address", ":8771"]
EXPOSE 8771
