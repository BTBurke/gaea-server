FROM ubuntu:14.04
MAINTAINER Bryan Burke <btburke@fastmail.com>

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y ca-certificates

EXPOSE 8080
EXPOSE 5432
EXPOSE 6379

COPY ./gaea-server /usr/local/bin/gaea-server
COPY ./entrypoint.sh /entrypoint.sh
ENV GIN_MODE release

ENTRYPOINT  ["/usr/local/bin/gaea-server"]
