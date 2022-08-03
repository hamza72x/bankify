include app.env
export

DB_URL="postgresql://${DB_USERNAME}:${DB_PASSWORD}@localhost:5432/bankify?sslmode=disable"
PWD=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))

up:
	docker-compose rm -f && \
	docker-compose build && \
	docker-compose --env-file ./app.env up

# use this while developing / testing with `air`
# `air` => live reloading for golang
db_dev:
	docker run -v "$(PWD)/tmp/vol_postgresql:/var/lib/postgresql/data" \
	--name ct_postgres --rm -p 5432:5432 \
	-e POSTGRES_USER="$(DB_USERNAME)" -e POSTGRES_DB="$(DB_NAME)" \
	-e POSTGRES_PASSWORD="$(DB_PASSWORD)" \
	postgres:latest

sqlc:
	sqlc generate

test:
	docker-compose --env-file ./app.env run service_test

# create migration file: 
# migrate create -ext sql -dir db/migration -seq $name 
# name example: init_schema

migrateup:
	migrate -path=./db/migration -database $(DB_URL) -verbose up

migrateup1:
	migrate -path=./db/migration -database $(DB_URL) -verbose up 1

migratedown:
	migrate -path=./db/migration -database $(DB_URL) -verbose down

migratedown1:
	migrate -path=./db/migration -database $(DB_URL) -verbose down 1

.PHONY: up dev_db test migrateup migrateup1 migratedown migratedown1 sqlc
