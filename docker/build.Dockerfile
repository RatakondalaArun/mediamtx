#################################################################
# Stage 1: Build
#################################################################
FROM golang:1.25-alpine3.22 AS builder

RUN apk add --no-cache git make
WORKDIR /s
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go generate ./...
RUN go build -o mediamtx

#################################################################
# Stage 2: Runtime
#################################################################
FROM alpine:3.22

RUN apk add --no-cache ffmpeg

COPY --from=builder /s/mediamtx /mediamtx
COPY --from=builder /s/mediamtx.yml /mediamtx.yml
COPY --from=builder /s/LICENSE /LICENSE

ENTRYPOINT [ "/mediamtx" ]
