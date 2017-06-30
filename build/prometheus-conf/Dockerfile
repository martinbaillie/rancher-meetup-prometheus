FROM martinbaillie/rancher-confd-onbuild
MAINTAINER Martin Baillie <martin.t.baillie@gmail.com>
ENV CONFD_WATCH true
VOLUME /etc/prometheus-conf
VOLUME /usr/local/bin
COPY rules/ /etc/prometheus-conf/rules
