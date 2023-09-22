FROM nlnetlabs/routinator:v0.12.1
USER root
COPY *.json *.sh *.crt routinator.conf /opt/
COPY *.py /home/routinator/
COPY --from=trajano/alpine-libfaketime  /faketime.so /lib/faketime.so
ENV LD_PRELOAD=/lib/faketime.so
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk update 
RUN apk add bash
RUN apk add --update --no-cache curl jq py3-configobj py3-setuptools python3 python3-dev 
RUN chmod -R 777 /opt/entrypoint.sh 
USER routinator
WORKDIR /home/routinator
ENTRYPOINT ["/opt/entrypoint.sh"]