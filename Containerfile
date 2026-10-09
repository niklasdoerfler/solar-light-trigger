FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG VERSION=dev
ARG BUILD_TIME=-

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=${TARGETVARIANT#v} \
    go build -trimpath \
    -ldflags "-s -w -X 'main.BuildVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'" \
    -o /out/solar-light-trigger .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="solar-light-trigger" \
      org.opencontainers.image.description="Evaluates solar radiation values via MQTT to determine day and night" \
      org.opencontainers.image.licenses="GPL-3.0"

WORKDIR /app
COPY --from=build /out/solar-light-trigger /app/solar-light-trigger

# Optional: mount a config.yaml to /config/config.yaml
VOLUME ["/config"]

ENTRYPOINT ["/app/solar-light-trigger"]
