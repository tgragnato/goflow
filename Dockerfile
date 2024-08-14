FROM golang:alpine3.20 as builder
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download && go build .

FROM alpine:3.20
WORKDIR /tmp
COPY --from=builder /workspace/goflow /usr/bin/
ENTRYPOINT ["/usr/bin/goflow"]
LABEL org.opencontainers.image.source=https://github.com/tgragnato/goflow
