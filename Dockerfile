FROM golang:1.25-alpine AS builder

LABEL org.opencontainers.image.source="https://github.com/armandwipangestu/gis-api"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server ./main.go

# Final stage
FROM alpine:3.18

ARG APP_ENV=development
ENV APP_ENV=$APP_ENV

ARG APP_VERSION=v1.0.0
ENV APP_VERSION=$APP_VERSION

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/server /server

WORKDIR /app

EXPOSE 3090

CMD ["/server"]