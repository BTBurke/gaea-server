FROM ubuntu:14.04
MAINTAINER Bryan Burke <btburke@fastmail.com>

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y ca-certificates

EXPOSE 8080
RUN mkdir -p /code
RUN mkdir -p /code/files
WORKDIR /code

ADD ./gaea-server /code/gaea-server
ENV GIN_MODE release

ENTRYPOINT  ["gaea-server"]
