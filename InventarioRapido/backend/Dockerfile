# Use a imagem oficial do Go
FROM golang:1.22-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory-api main.go

# Usar imagem mínima para produção
FROM alpine:latest

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -S -u 1001 -G appgroup appuser

WORKDIR /root/

# Copiar binário da etapa anterior
COPY --from=builder /app/inventory-api .
COPY --from=builder /app/.env .

# Mudar proprietário dos arquivos
RUN chown -R appuser:appgroup /root/

# Usar usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Comando para executar
CMD ["./inventory-api"]