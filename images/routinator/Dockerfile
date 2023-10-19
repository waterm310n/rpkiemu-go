FROM alpine:3.17.2 AS builder
FROM nlnetlabs/routinator:v0.12.1
USER root
RUN apk add bash 
COPY *.json *.sh *.crt routinator.conf /opt/ 
RUN chmod -R 777 /opt/entrypoint.sh 
WORKDIR /root
ENTRYPOINT [""]
CMD ["/bin/sleep","2000000"]