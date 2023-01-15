FROM golang:1.19.3-alpine as base
WORKDIR /app

FROM base as builder
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go generate ./...
RUN go build -o toiyeugo server.go

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /app/toiyeugo toiyeugo
ENTRYPOINT ["./toiyeugo"]