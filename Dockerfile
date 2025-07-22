FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/build/app /app/cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/build/app /app/build/app
EXPOSE 3000
CMD [ "./build/app" ]