#!/bin/bash
set -e

APP_USER_PASSWORD=$(cat /run/secrets/postgres_password)

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO \$$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'app_user') THEN
            EXECUTE format('CREATE ROLE app_user WITH LOGIN PASSWORD %L', '$APP_USER_PASSWORD');
            ALTER ROLE app_user WITH LOGIN;
        END IF;
    END
    \$$;

    GRANT CONNECT ON DATABASE "$POSTGRES_DB" TO app_user;
    GRANT USAGE ON SCHEMA public TO app_user;
EOSQL