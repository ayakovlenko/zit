FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /zit ./cmd/zit

FROM alpine:3
RUN apk add --no-cache git
COPY --from=builder /zit /usr/local/bin/zit
