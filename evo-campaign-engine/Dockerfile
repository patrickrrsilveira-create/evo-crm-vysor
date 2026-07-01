FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/campaign-engine ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/campaign-engine ./campaign-engine
COPY --from=builder /go/bin/migrate ./migrate
COPY migrations ./migrations
COPY entrypoint.sh ./entrypoint.sh

RUN sed -i 's/\r$//' ./entrypoint.sh
RUN chmod +x ./campaign-engine ./migrate ./entrypoint.sh

EXPOSE 9090

CMD ["./entrypoint.sh"]
