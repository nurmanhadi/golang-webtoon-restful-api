FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base libwebp-dev

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/build/main /app/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/build/main /app/build/main
EXPOSE 3000
CMD [ "./build/main" ]