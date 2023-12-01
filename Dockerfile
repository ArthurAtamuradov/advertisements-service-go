# Dockerfile
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app /app

EXPOSE 8080

CMD ["./app"]
