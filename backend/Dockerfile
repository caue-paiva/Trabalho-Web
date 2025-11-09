# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS build

ENV GOTOOLCHAIN=auto \
    GOPROXY=https://proxy.golang.org,direct

# Ensure CA bundle exists for TLS during go mod download
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM scratch
WORKDIR /app
# Copy CA bundle if your app makes HTTPS calls
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/server ./server
COPY --from=build /src/configs ./configs
EXPOSE 8080
ENTRYPOINT ["./server"]
