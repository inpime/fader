FROM ubuntu:14.04

MAINTAINER gebv <egor.workroom@gmail.com>

RUN apt-get update && apt-get -y install curl netcat unzip

EXPOSE 3322

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY fader /app/fader
RUN chmod +x /app/fader

RUN mkdir -p /app && mkdir -p /app/logs

ENTRYPOINT ["/entrypoint.sh"]