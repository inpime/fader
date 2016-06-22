FROM alpine:3.3

MAINTAINER gebv <egor.workroom@gmail.com>

RUN mkdir -p /app

RUN apk update && apk upgrade \
 && apk update && apk add unzip \
 && apk --no-cache --no-progress add ca-certificates \
 && wget https://s3.eu-central-1.amazonaws.com/releases.fader.inpime.com/fader.go1.6.linux_amd64.latest.zip \
 && unzip fader.go1.6.linux_amd64.latest.zip -d /app \
 && chmod a+x /app/* \
 && chmod -R 777 /app/* \
 && rm -rf /var/cache/apk/*

EXPOSE 3322

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

WORKDIR /app
ENTRYPOINT ["fader"]
# ENTRYPOINT ["/entrypoint.sh"]
# ENTRYPOINT ["/bin/sh", "/app/fader"]
