FROM debian:13.3-slim

ENV SHELL=/bin/bash

RUN apt-get update

# install dependencies
RUN apt-get install -y cron lftp netcat-traditional zip lighttpd curl jq
RUN curl -Ls https://github.com/mikefarah/yq/releases/download/v4.53.6/yq_linux_amd64 -o /usr/bin/yq
RUN chmod +x /usr/bin/yq

# setup crontab
COPY ./etc/crontab /etc/
RUN chmod 600 /etc/crontab
RUN chown root:root /etc/crontab

# setup dashboard configuration
COPY ./etc/lighttpd.conf /etc/lighttpd/lighttpd.conf

# setup dashboard generator script
COPY ./bin/* /usr/local/bin/
RUN chmod +x /usr/local/bin/*

# setup dashboard templates
RUN mkdir -p /var/lib/backup_3ds/dashboard/templates
RUN mkdir -p /var/lib/backup_3ds/dashboard/static
COPY dashboard/templates/* /var/lib/backup_3ds/dashboard/templates/
COPY dashboard/static/* /var/lib/backup_3ds/dashboard/static/

EXPOSE 80

WORKDIR /usr/local/bin
CMD [ "./entrypoint.sh" ]
