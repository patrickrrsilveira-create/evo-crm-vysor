# -----------------------------------------------------
# Builder Stage: Compila o binário Go otimizado
# -----------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instala dependências de SO necessárias para o build
RUN apk add --no-cache gcc musl-dev

# Baixa as dependências do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático (-w -s remove símbolos de debug para reduzir tamanho)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o evo-swarm-engine ./cmd/server/main.go

# -----------------------------------------------------
# Production Stage: Imagem final super leve (Scratch ou Alpine)
# -----------------------------------------------------
FROM alpine:latest

WORKDIR /root/

# Instala certificados SSL e tzdata (para timezones)
RUN apk --no-cache add ca-certificates tzdata

# Copia o binário do stage anterior
COPY --from=builder /app/evo-swarm-engine .

# Copia as variáveis de ambiente (Em produção, idealmente injetadas via Secret, mas copiamos o default)
COPY .env .env

# Expõe a porta do Fiber
EXPOSE 8001

# Comando de inicialização
CMD ["./evo-swarm-engine"]
