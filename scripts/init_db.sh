#!/usr/bin/env bash
set -x
set -eo pipefail

export DOCKER_DEFAULT_PLATFORM=linux/amd64

if ! [ -x "$(command -v psql)" ]; then
    echo >&2 "Error: psql is not installed"
    exit 1
fi

DB_USER=${POSTGRES_USER:=postgres}
DB_PASSWORD=${POSTGRES_PASSWORD:=password}
DB_NAME=${POSTGRES_DB:=muttr}
DB_PORT=${POSTGRES_PORT:=5432}

if docker ps | grep postgres; then
    continue
elif docker ps -a | grep postgres; then
    docker start postgres
else
    docker run \
        --name postgres \
        -e POSTGRES_USER=${DB_USER} \
        -e POSTGRES_PASSWORD=${DB_PASSWORD} \
        -e POSTGRES_DB=${DB_NAME} \
        -e POSTGRES_PORT=${DB_PORT} \
        -p "${DB_PORT}":5432 \
        -d postgres \
        postgres -N 1000
fi

while true ; do
    if [ -x "$(command docker exec postgres psql -U postgres -c \"\\q\")"]; then
        DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}
        break
    else
        >&2 echo "Postgres is still unavailable - sleeping"
        sleep 1
    fi
done

for migration in $(ls ./internal/storage/migrations); do
    if [[ $migration == *"_up.sql" ]]; then
        docker exec postgres psql -U postgres -c "$(cat ./internal/storage/migrations/${migration})";
    fi
done