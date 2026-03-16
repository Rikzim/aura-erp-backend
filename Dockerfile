FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/backend-go ./main.go

FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/backend-go /app/backend-go

ENV PORT=5000 \
    GIN_MODE=release

EXPOSE 5000
USER app

ENTRYPOINT ["/app/backend-go"]
