FROM douyin-cloud-cn-beijing.cr.volces.com/cloud-public/builder:alpine-3.13


run sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories 

RUN apk update \
  &&  apk upgrade  \
  &&  apk add bash \
  && apk add curl

WORKDIR /opt/application

RUN ./build.sh
COPY output/main /opt/application/main
COPY run.sh /opt/application/

USER root
CMD cd /opt/application && ./main