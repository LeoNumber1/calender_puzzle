FROM golang:1.16.5 as builder
COPY . /go/src/puzzle
WORKDIR /go/src/puzzle
RUN make build-linux

FROM debian:bullseye
RUN mkdir -p /app
COPY --from=builder /go/src/puzzle/bin /app

WORKDIR     /app
ENTRYPOINT  ["/app/puzzle-linux"]
