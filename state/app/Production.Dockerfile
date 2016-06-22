FROM ubuntu:14.04

MAINTAINER gebv <egor.workroom@gmail.com>

RUN mkdir -p /app

RUN apt-get update && apt-get -y install curl netcat unzip

RUN curl -O https://s3.eu-central-1.amazonaws.com/releases.fader.inpime.com/fader.go1.6.linux_amd64.latest.zip \
 && unzip fader.go1.6.linux_amd64.latest.zip -d /app \
 && chmod +x /app/fader

EXPOSE 3322

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
