FROM alpine:3.5
MAINTAINER community@apstra.com

COPY telegraf /bin/telegraf
RUN chmod +x /bin/telegraf

RUN ln -s /bin/telegraf /usr/bin/telegraf
