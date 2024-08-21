FROM   bitnami/postgresql:latest
EXPOSE 5432
ENV    POSTGRESQL_PASSWORD=password123 \
       POSTGRESQL_DATABASE=readcommend \
       DB_URL=postgres://postgres:password123@localhost:5432/readcommend?sslmode=disable
COPY   . ./docker-entrypoint-initdb.d