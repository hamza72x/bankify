FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./

COPY api ./api
COPY config ./config
COPY db ./db
COPY util ./util

RUN go build -o /app/runner

# COPY ../build/bin/golang-migrate/migrate /usr/bin/migrate
# RUN migrate -path db/migration -database "$(DB_URL)" -verbose up
# RUN apk add wkhtmltopdf

# run the app
FROM golang:1.18

WORKDIR /app

COPY --from=builder /app/runner .
COPY app.env ./
COPY db/migration ./db/migration

CMD ["/app/runner"]
