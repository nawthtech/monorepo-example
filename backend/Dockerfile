FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# بناء الخادم فقط
RUN CGO_ENABLED=0 GOOS=linux go build -o server-app ./cmd/server

# مرحلة التشغيل
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server-app .

EXPOSE 8080

CMD ["./server-app"]