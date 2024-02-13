From nlnetlabs/krill:v0.12.1 
USER root
COPY *.sh /var/krill/data/

COPY --from=trajano/alpine-libfaketime  /faketime.so /lib/faketime.so
ENV LD_PRELOAD=/lib/faketime.so