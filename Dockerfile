# Etapa de build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o load-tester .

# Imagem final
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/load-tester .

ENTRYPOINT ["./load-tester"]
