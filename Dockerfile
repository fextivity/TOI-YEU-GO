FROM golang:1.19.3-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o toiyeugo

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=builder /app/toiyeugo toiyeugo
ENTRYPOINT ["./toiyeugo"]