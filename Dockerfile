# Build
FROM golang:1.17.0-alpine AS builder

ENV GO111MODULE=on
ENV GOPATH=/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o /main cmd/main.go

# Deploy
FROM alpine

COPY --from=builder main .

CMD ["./main"]