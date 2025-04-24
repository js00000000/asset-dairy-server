FROM golang:1.24.2-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy

FROM golang:1.24.2-alpine
WORKDIR /app
COPY --from=builder /app /app
RUN apk add --no-cache git
EXPOSE 3000
CMD ["go", "run", "main.go"]
