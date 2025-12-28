#################################################################
# Stage 1: Build binaries for all platforms
#################################################################
ARG BASE_IMAGE=golang:1.25-alpine3.22

FROM ${BASE_IMAGE} AS build-base
RUN apk add --no-cache zip make git tar
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ENV CGO_ENABLED=0
RUN rm -rf tmp binaries
RUN mkdir tmp binaries
RUN cp mediamtx.yml LICENSE tmp/
RUN go generate ./...

FROM build-base AS build-linux-amd64
ENV GOOS=linux GOARCH=amd64
RUN go build -tags enableUpgrade -o "tmp/mediamtx"
RUN tar -C tmp -czf "binaries/mediamtx_$(cat internal/core/VERSION)_linux_amd64.tar.gz" --owner=0 --group=0 "mediamtx" mediamtx.yml LICENSE

FROM build-base AS build-linux-armv6
ENV GOOS=linux GOARCH=arm GOARM=6
RUN go build -tags enableUpgrade -o "tmp/mediamtx"
RUN tar -C tmp -czf "binaries/mediamtx_$(cat internal/core/VERSION)_linux_armv6.tar.gz" --owner=0 --group=0 "mediamtx" mediamtx.yml LICENSE

FROM build-base AS build-linux-armv7
ENV GOOS=linux GOARCH=arm GOARM=7
RUN go build -tags enableUpgrade -o "tmp/mediamtx"
RUN tar -C tmp -czf "binaries/mediamtx_$(cat internal/core/VERSION)_linux_armv7.tar.gz" --owner=0 --group=0 "mediamtx" mediamtx.yml LICENSE

FROM build-base AS build-linux-arm64
ENV GOOS=linux GOARCH=arm64
RUN go build -tags enableUpgrade -o "tmp/mediamtx"
RUN tar -C tmp -czf "binaries/mediamtx_$(cat internal/core/VERSION)_linux_arm64.tar.gz" --owner=0 --group=0 "mediamtx" mediamtx.yml LICENSE

#################################################################
# Stage 2: Collect all binaries
#################################################################
FROM ${BASE_IMAGE} as binaries-collector
RUN apk add --no-cache alpine-sdk
COPY --from=build-linux-amd64 /s/binaries /binaries/linux/amd64
COPY --from=build-linux-armv6 /s/binaries /binaries/linux/arm/v6
COPY --from=build-linux-armv7 /s/binaries /binaries/linux/arm/v7
COPY --from=build-linux-arm64 /s/binaries /binaries/linux/arm64

#################################################################
# Stage 3: Final runtime image with ffmpeg
#################################################################
FROM alpine:3.22

RUN apk add --no-cache ffmpeg

ARG TARGETPLATFORM
COPY --from=binaries-collector /binaries/${TARGETPLATFORM} /

ENTRYPOINT [ "/mediamtx" ]
