#! /bin/bash

source .env

sqlc generate

atlas schema apply \
    --url "postgres://$PG_USER:$PG_PASSWORD@localhost:5432/$PG_DBNAME?sslmode=disable" \
    --dev-url "docker://postgres" \
    --to "file://db/query/schema.sql"
