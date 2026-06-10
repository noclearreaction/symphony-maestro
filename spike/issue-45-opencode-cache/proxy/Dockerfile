FROM golang:1.26 AS builder
WORKDIR /build
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o proxy .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/proxy /proxy
EXPOSE 8080
ENTRYPOINT ["/proxy"]
