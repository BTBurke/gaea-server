FROM ubuntu:14.04
MAINTAINER Bryan Burke <btburke@fastmail.com>

COPY ./gaea-server /usr/local/bin/gaea-server
COPY ./entrypoint.sh /entrypoint.sh
ENV GIN_MODE release

ENTRYPOINT  ["/entrypoint.sh"]
