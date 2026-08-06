# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/erpagendamentos ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -g 10001 appgroup \
    && adduser -D -u 10001 -G appgroup appuser
WORKDIR /app
COPY --from=build /out/erpagendamentos /app/erpagendamentos
COPY sql /app/sql
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/erpagendamentos"]
