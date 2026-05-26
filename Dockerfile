FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bookshelf ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/bookshelf .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
RUN mkdir -p /data
EXPOSE 8080
CMD ["./bookshelf"]
