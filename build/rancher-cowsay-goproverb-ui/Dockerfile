FROM alpine
MAINTAINER Martin Baillie <martin.t.baillie@gmail.com>

RUN apk add --no-cache openssh-client git tar curl && \
        curl -k --silent --show-error --fail --location --connect-timeout 360 \
        --header "Accept: application/tar+gzip, application/x-gzip, application/octet-stream" -o - \
        "https://caddyserver.com/download/linux/amd64?plugins=http.prometheus,http.cors" \
        | tar --no-same-owner -C /usr/bin/ -xz caddy \
        && chmod 0755 /usr/bin/caddy

WORKDIR /srv

COPY Caddyfile /etc/Caddyfile
COPY index.html prom.js cow.png /srv/

ENTRYPOINT ["/usr/bin/caddy"]
CMD ["--conf", "/etc/Caddyfile", "--log", "stdout"]
