FROM cgr.dev/chainguard/go:latest as builder
ENV CGO_ENABLED=0
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download && go build .

FROM ghcr.io/anchore/syft:latest AS sbomgen
COPY --from=builder /workspace/goflow /usr/bin/goflow
RUN ["/syft", "--output", "spdx-json=/goflow.spdx.json", "/usr/bin/goflow"]

FROM cgr.dev/chainguard/static:latest
WORKDIR /tmp
COPY --from=builder /workspace/goflow /usr/bin/
COPY --from=sbomgen /goflow.spdx.json /var/lib/db/sbom/goflow.spdx.json
ENTRYPOINT ["/usr/bin/goflow"]
LABEL org.opencontainers.image.title="goflow"
LABEL org.opencontainers.image.description="The high-scalability sFlow/NetFlow/IPFIX collector used internally at Cloudflare"
LABEL org.opencontainers.image.url="https://tgragnato.it/goflow/"
LABEL org.opencontainers.image.source="https://tgragnato.it/goflow/"
LABEL license="BSD-3-Clause"
LABEL io.containers.autoupdate=registry