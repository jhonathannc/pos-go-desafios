# build
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o loadtester main.go

# imagem final
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/loadtester .

ENTRYPOINT ["./loadtester"]
