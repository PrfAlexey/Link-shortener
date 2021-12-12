FROM golang:1.15 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./cmd/main.go


FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start && \
  psql --command "CREATE USER admin WITH SUPERUSER PASSWORD '4444';" && \
  createdb -O admin link-shortener && \
  /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

EXPOSE 5000
ENV PGPASSWORD 4444


CMD service postgresql start &&  psql -h localhost -d link-shortener -U admin -p 5432 -a -q -f schema.sql  && ./main