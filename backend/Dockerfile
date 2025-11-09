# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build

ENV GOTOOLCHAIN=auto \
    GOPROXY=https://proxy.golang.org,direct

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/server ./server
COPY --from=build /src/configs ./configs
EXPOSE 8080
ENTRYPOINT ["./server"]
