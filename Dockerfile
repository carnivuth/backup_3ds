FROM debian:13.3-slim

RUN apt update

# install cron and crontab
RUN apt install -y cron ncftp netcat-traditional

COPY crontab /etc/
# set permissions
RUN chmod 600 /etc/crontab
RUN chown root:root /etc/crontab

# add project files
WORKDIR /usr/local/bin
COPY backup_3ds.sh .
COPY start.sh .

RUN chmod +x backup_3ds.sh start.sh

CMD [ "./start.sh" ]
