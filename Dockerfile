FROM ubuntu:14.04
MAINTAINER Bryan Burke <btburke@fastmail.com>

COPY ./gaea-server /usr/local/bin/gaea-server
ENV GIN_MODE release

ENTRYPOINT  ["/usr/local/bin/gaea-server"]
