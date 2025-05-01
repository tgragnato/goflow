FROM cgr.dev/chainguard/go:latest as builder
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download && go build .

FROM cgr.dev/chainguard/static:latest
WORKDIR /tmp
COPY --from=builder /workspace/goflow /usr/bin/
ENTRYPOINT ["/usr/bin/goflow"]
LABEL org.opencontainers.image.title="goflow"
LABEL org.opencontainers.image.description="The high-scalability sFlow/NetFlow/IPFIX collector used internally at Cloudflare"
LABEL org.opencontainers.image.url="https://tgragnato.it/goflow/"
LABEL org.opencontainers.image.source="https://tgragnato.it/goflow/"
LABEL license="BSD-3-Clause"
LABEL io.containers.autoupdate=registry