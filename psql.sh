#!/bin/bash
# Connects to database for development purposes
. vars.sh
docker-compose exec postgresql psql $DB_URL
