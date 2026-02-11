FROM debian:13.3-slim

RUN apt update

# install cron and crontab
RUN apt install -y cron ncftp netcat-traditional zip lighttpd

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
COPY dashboard/style/* /var/lib/backup_3ds/dashboard/static/

WORKDIR /usr/local/bin
CMD [ "./entrypoint.sh" ]
