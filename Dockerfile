FROM golang:1.19.5-alpine as base
WORKDIR /app

FROM base as builder
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o toiyeugo

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /app/toiyeugo toiyeugo
ENTRYPOINT ["./toiyeugo"]