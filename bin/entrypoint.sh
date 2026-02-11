#! /bin/bash
 printenv | grep -v no_proxy >> /etc/environment
 echo "starting cron now"
 /usr/sbin/cron -f
