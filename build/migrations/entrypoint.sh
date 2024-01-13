#!/bin/bash

DBSTRING="host=$OHA_POSTGRESQL_HOST port=$OHA_POSTGRESQL_PORT user=$OHA_POSTGRESQL_USERNAME password=$OHA_POSTGRESQL_PASSWORD dbname=$OHA_POSTGRESQL_DBNAME sslmode=disable"

goose postgres "$DBSTRING" up