# Construir el binario de la aplicación
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Crear la imagen final para la ejecución
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .
COPY --from=builder /app/dataBase /root/dataBase
COPY --from=builder /app/handlers /root/handlers
EXPOSE 8080
CMD ["./main"]
