FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/strava-kudos ./strava-kudos.go

FROM scratch
COPY --from=builder /app/strava-kudos /app/strava-kudos
WORKDIR /app
ENTRYPOINT ["/app/strava-kudos"]