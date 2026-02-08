#!/bin/bash

log_inner() {
        echo "[${BASH_SOURCE##*/}:${1^^}] ${FUNCNAME[2]}@${BASH_LINENO[1]}: ${*:2}"
}

# check on env vars
if [[ -z $FTPD_3DS_ADDRESS ]]; then log_inner error "Set server address trough FTPD_3DS_ADDRESS env var"; exit 1; fi
if [[ -z $FTPD_3DS_PORT ]]; then FTPD_3DS_PORT=21; exit 1; fi

# directories
if [[ -z $BASE_DIR ]];then BASE_DIR=/var/lib/3ds_backup; fi
if [[ ! -d "$BASE_DIR" ]]; then mkdir -p "$BASE_DIR"; fi

if [[ -z $BACKUP_DEST ]];then BACKUP_DEST="$BASE_DIR/backups"; fi
if [[ ! -d "$BACKUP_DEST" ]]; then mkdir -p "$BACKUP_DEST"; fi

# stat file to share information between cron instances of the script
if [[ -z $STAT_FILE ]];then STAT_FILE="/tmp/status"; fi
if [[ ! -f "$STAT_FILE" ]]; then echo 1 > "$STAT_FILE"; fi

# 3ds directory to backup
if [[ -z $BACKUP_SRC ]];then BACKUP_SRC="/"; fi

function backup(){

  # check if 3ds ftp server is up
  if ! nc -z -w1 "$FTPD_3DS_ADDRESS" "$FTPD_3DS_PORT"; then log_inner info "3ds is not listening for ftp connections"; exit 0; fi

  # check if 3ds backup has been made
  if [[ $(cat "$STAT_FILE") == "0" ]]; then log_inner info "backup has already been made at $(stat -c '%y' "$STAT_FILE")"; exit 0; fi

  # check if 3ds backup is running
  if [[ $(cat "$STAT_FILE") == "2" ]]; then log_inner info "backup is running since $(stat -c '%y' "$STAT_FILE")"; exit 0; fi

  log_inner info "starting backup of 3DS address at $FTPD_3DS_ADDRESS:$FTPD_3DS_PORT"
  # setting status on status file to avoid running multiple backup jobs
  echo 2 > "$STAT_FILE"

  timestamp="$(date +%s)"
  IFS=';' read -ra dirs <<< "$BACKUP_SRC"
  for dir in ${dirs[@]}; do

    dirname="$(echo "$dir" | sed 's/\//-/g')"
    mkdir "$BACKUP_DEST/${timestamp}_$dirname"
    log_inner info "creating backup $BACKUP_DEST/${timestamp}_$dir of $dirname"

    # check for error codes and print error otherwise
    if ncftpget -R -v -u "$FTPD_3DS_USERNAME" -p "$FTPD_3DS_PASSWORD" -P "$FTPD_3DS_PORT" "$FTPD_3DS_ADDRESS" "$BACKUP_DEST/${timestamp}_${dirname}" "$dir"; then
      tar -czf "$BACKUP_DEST/${timestamp}_${dirname}.tar.gz" -C "$BACKUP_DEST" "${timestamp}_${dirname}" || log_inner error "error archiving ${BACKUP_DEST}/${timestamp}_${dirname}"
      rm -r "$BACKUP_DEST/${timestamp}_$dirname"
      log_inner info "done backup $BACKUP_DEST/${timestamp}_$dirname"
    else
      log_inner error "error in downloading $dir from $FTPD_3DS_ADDRESS"
    fi
  done

  # setting status file to "backup done" to avoid consequent backups
  echo 0 > "$STAT_FILE"
}

function reset(){

  # check if 3ds backup is running
  if [[ $(cat "$STAT_FILE") == "2" ]]; then log_inner info "backup is running since $(stat -c '%y' "$STAT_FILE"), avoid resetting"; exit 0; fi

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
    log_inner error "usage: $0 backup/reset"
    ;;
esac
