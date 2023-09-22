# Warning: there is only a 'latest' tag, there is no version specific tag that
# we can freeze at.
FROM alpine:3.17.2 AS builder
FROM vimagick/rsyncd
USER root
RUN apk add bash 
COPY rsyncd.conf /etc/rsyncd.conf
COPY *.sh /opt/
USER root
RUN chmod -R 777 /opt/entrypoint.sh
ENTRYPOINT [ "/opt/entrypoint.sh" ]