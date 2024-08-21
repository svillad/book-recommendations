#!/bin/bash
# Creates the database tables and seeds them with sample data
psql -a -q -f /docker-entrypoint-initdb.d/migrate.psql $DB_URL
