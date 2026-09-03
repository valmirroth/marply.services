# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/bh-service .

# --- runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build /out/bh-service ./bh-service

ENV HTTP_PORT=:8093
EXPOSE 8093

ENTRYPOINT ["./bh-service"]
