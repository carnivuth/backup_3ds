#!/bin/bash

# check on env vars
if [[ -z $FTPD_3DS_ADDRESS ]]; then echo "Set server address trough FTPD_3DS_ADDRESS env var"; exit 1; fi
if [[ -z $FTPD_3DS_PORT ]]; then echo "Set server port trough FTPD_3DS_PORT env var"; exit 1; fi

# directories
BASE_DIR=/var/lib/3ds_backup; if [[ ! -d "$BASE_DIR" ]]; then mkdir -p "$BASE_DIR"; fi
BACKUP_DEST="$BASE_DIR/backup"; if [[ ! -d "$BACKUP_DEST" ]]; then mkdir -p "$BACKUP_DEST"; fi
STAT_FILE="$BASE_DIR/status"; if [[ ! -f "$STAT_FILE" ]]; then echo 1 > "$STAT_FILE"; fi

function backup(){

  # check if 3ds ftp server is up
  if ! nc -z -w1 "$FTPD_3DS_ADDRESS" "$FTPD_3DS_PORT"; then echo "3ds is not listening for ftp connections"; exit 0; fi

  # check if 3ds backup is been made
  if [[ $(cat "$STAT_FILE") == "0" ]]; then echo "backup has already been made at $(stat -c '%y' "$STAT_FILE")"; exit 0; fi

  # check if 3ds backup is running
  if [[ $(cat "$STAT_FILE") == "2" ]]; then echo "backup is running since $(stat -c '%y' "$STAT_FILE")"; exit 0; fi

  echo "connecting to $FTPD_3DS_ADDRESS $FTPD_3DS_PORT"
  echo 2 > "$STAT_FILE"
  timestamp="$(date +%s)"
  mkdir $BACKUP_DEST/$timestamp
  ncftpget -d stdout -R -P "$FTPD_3DS_PORT" "$FTPD_3DS_ADDRESS" "$BACKUP_DEST/$timestamp" /
  tar -czvf  "$BACKUP_DEST/$timestamp.tar.gz" "$BACKUP_DEST/$timestamp"
  echo "done backup"
  echo 0 > "$STAT_FILE"

}

function reset(){

  # check if 3ds backup is running
  if [[ $(cat "$STAT_FILE") == "2" ]]; then echo "backup is running since $(stat -c '%y' "$STAT_FILE"), avoid resetting"; exit 0; fi

  echo 1 > "$STAT_FILE"
}

case "$1" in
  "backup")
    backup
    ;;
  "reset")
    reset
    ;;
  *)
    echo "usage: $0 backup/reset"
    ;;
esac
