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
COPY 3ds_backup.sh .
COPY start.sh .

RUN chmod +x 3ds_backup.sh start.sh

CMD [ "./start.sh" ]
