# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Instala dependências do sistema
RUN apk add --no-cache git

# Instala golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copia os arquivos de dependência
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Debug: Verifica se o binário foi criado
RUN ls -la /app/main

# Final stage
FROM alpine:latest

WORKDIR /app

# Instala ca-certificates para HTTPS requests
RUN apk --no-cache add ca-certificates

# Copia o binário compilado do stage anterior
COPY --from=builder /app/main ./main

# Copia o migrate do builder
COPY --from=builder /go/bin/migrate ./migrate

# Copia as migrations
COPY migrations ./migrations

# Copia o script de inicialização
COPY entrypoint.sh ./entrypoint.sh

# Debug: Lista arquivos no diretório
RUN ls -la /app/

# Torna os binários e script executáveis
RUN chmod +x ./main ./migrate ./entrypoint.sh

# Debug: Verifica se os arquivos existem e suas permissões
RUN ls -la ./main ./migrate ./entrypoint.sh

# Expõe a porta da aplicação
EXPOSE 5555

# Executa o script de inicialização
CMD ["./entrypoint.sh"] 