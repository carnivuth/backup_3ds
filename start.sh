#! /bin/bash
 printenv | grep -v no_proxy >> /etc/environment
 /usr/sbin/cron -f
