FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/strava-kudos ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/strava-kudos /app/strava-kudos
COPY .env /app/.env
COPY .strava-auth-token /app/.strava-auth-token
WORKDIR /app
ENTRYPOINT ["/app/strava-kudos"]