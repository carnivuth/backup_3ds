#!/bin/bash

log_inner() {
        echo "[${BASH_SOURCE##*/}:${1^^}] ${FUNCNAME[2]}@${BASH_LINENO[1]}: ${*:2}"
}

# check if 3ds address variable is set, exit otherwise
if [[ -z $FTPD_3DS_ADDRESSES ]]; then log_inner error "Set server address trough FTPD_3DS_ADDRESSES env var"; exit 1; fi

# setting default port
if [[ -z $FTPD_3DS_PORT ]]; then FTPD_3DS_PORT=21; fi

# directories
if [[ -z $BASE_DIR ]];then BASE_DIR=/var/lib/3ds_backup; fi
if [[ ! -d "$BASE_DIR" ]]; then mkdir -p "$BASE_DIR"; fi

if [[ -z $BACKUP_DEST ]];then BACKUP_DEST="$BASE_DIR/backups"; fi
if [[ ! -d "$BACKUP_DEST" ]]; then mkdir -p "$BACKUP_DEST"; fi

# stat file to share information between cron instances of the script
if [[ -z $STAT_DIR ]];then STAT_DIR="/tmp/backup_3ds/status"; fi
if [[ ! -d "$STAT_DIR" ]]; then mkdir -p "$STAT_DIR"; fi

# 3ds directory to backup
if [[ -z $BACKUP_SRC ]];then BACKUP_SRC="/"; fi

function backup(){

  IFS=';' read -ra addresses <<< "$FTPD_3DS_ADDRESSES"
  for address in ${addresses[@]}; do

    # set status file variable
    STAT_FILE="${STAT_DIR}/${address}"
    if [[ ! -f "$STAT_FILE" ]]; then echo 1 > "$STAT_FILE"; fi

    # setting backup dir for the specific 3ds
    if [[ -z $HOST_DIR ]];then HOST_DIR="${BACKUP_DEST}/${address}"; fi
    if [[ ! -d "$HOST_DIR" ]]; then mkdir -p "$HOST_DIR"; fi

    # check if 3ds backup is not running
    if [[ $(cat "$STAT_FILE") != "2" ]]; then

      # check if 3ds backup has not been made
      if [[ $(cat "$STAT_FILE") != "0" ]]; then

        # check if 3ds ftp server is up
        if nc -z -w1 "$address" "$FTPD_3DS_PORT"; then
          # setting lock file to avoid running multiple backup jobs in parallel on the same 3DS
          echo 2 > "$STAT_FILE"
          log_inner info "starting backup of 3DS at $address:$FTPD_3DS_PORT"

          timestamp="$(date +%s)"

          IFS=';' read -ra dirs <<< "$BACKUP_SRC"
          for dir in ${dirs[@]}; do

            dirname="$(echo "$dir" | sed 's/\//-/g')"
            mkdir "$HOST_DIR/${timestamp}_$dirname"
            log_inner info "creating backup $HOST_DIR/${timestamp}_$dir of $dirname"

            # check for error codes and print error otherwise
            if ncftpget -R -v -u "$FTPD_3DS_USERNAME" -p "$FTPD_3DS_PASSWORD" -P "$FTPD_3DS_PORT" "$address" "$HOST_DIR/${timestamp}_${dirname}" "$dir"; then
              tar -czf "$HOST_DIR/${timestamp}_${dirname}.tar.gz" -C "$HOST_DIR" "${timestamp}_${dirname}" || log_inner error "error archiving ${HOST_DIR}/${timestamp}_${dirname}"
              rm -r "${HOST_DIR}/${timestamp}_$dirname"
              log_inner info "done backup ${HOST_DIR}/${timestamp}_$dirname"
            else
              log_inner error "error in downloading ${dir} from ${address}"
            fi
          done

          # setting status file to "backup done" to avoid consequent backups
          echo 0 > "$STAT_FILE"

        else
          log_inner info "3ds is not listening for ftp connections"
        fi

      else
        log_inner info "backup has already been made at $(stat -c '%y' "$STAT_FILE")";
      fi

    else
      log_inner info "backup is running since $(stat -c '%y' "$STAT_FILE")";
    fi

  done



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
