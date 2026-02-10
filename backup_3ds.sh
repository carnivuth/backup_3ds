#!/bin/bash

log_inner() {
        echo "[${BASH_SOURCE##*/}:${1^^}] ${FUNCNAME[2]}@${BASH_LINENO[1]}: ${*:2}"
}

# check if 3ds address variable is set, exit otherwise
if [[ -z $FTPD_3DS_ADDRESSES ]]; then log_inner error "Set server address trough FTPD_3DS_ADDRESSES env var"; exit 1; fi

# ports parsing, if no ports are setted in the env var fallback to 21 for all host
if [[ -z $FTPD_3DS_PORTS ]]; then FTPD_3DS_PORTS='21'; fi

# directories
if [[ -z $BASE_DIR ]];then BASE_DIR=/var/lib/backup_3ds; fi
if [[ ! -d "$BASE_DIR" ]]; then mkdir -p "$BASE_DIR"; fi

if [[ -z $BACKUP_DEST ]];then BACKUP_DEST="$BASE_DIR/backups"; fi
if [[ ! -d "$BACKUP_DEST" ]]; then mkdir -p "$BACKUP_DEST"; fi

# stat dir to share information between cron instances of the script, this act as a folder for lock files
if [[ -z $STAT_DIR ]];then STAT_DIR="/tmp/backup_3ds/status"; fi
if [[ ! -d "$STAT_DIR" ]]; then mkdir -p "$STAT_DIR"; fi

# parsing the 3DS directory to backup
if [[ -z $BACKUP_DIRS ]];then BACKUP_DIRS="/"; fi

# backup a single 3ds, parameters address port username password
function backup(){

  if [[ -z $1 ]];then log_inner error "pass address as parameter"; return 4; else address="$1"; fi
  if [[ -z $2 ]];then port=21; else port="$2"; fi
  username=$3
  password=$4

  # set status file variable and initialize the file if it does not exist
  stat_file="${STAT_DIR}/${address}"
  if [[ ! -f "$stat_file" ]]; then echo 1 > "$stat_file"; fi

  # setting backup dir for the specific 3ds
  host_dir="${BACKUP_DEST}/${address}"
  if [[ ! -d "$host_dir" ]]; then mkdir -p "$host_dir"; fi

  # check if 3ds backup is running
  if [[ $(cat "$stat_file") == "2" ]]; then log_inner info "backup is running since $(stat -c '%y' "$stat_file")"; return 2; fi

  # check if 3ds backup has been made
  if [[ $(cat "$stat_file") == "0" ]]; then log_inner info "backup has already been made at $(stat -c '%y' "$stat_file")"; return 1; fi

  # check if 3ds ftp server is up
  if ! nc -z -w1 "$address" "$port"; then log_inner info "3ds at ${address}:${port} is not listening for ftp connections"; return 3; fi

  # setting lock file to avoid running multiple backup jobs in parallel on the same 3DS
  echo 2 > "$stat_file"
  log_inner info "starting 3DS backup of $address:$port"

  timestamp="$(date +%s)"

  # loop all dirs to backup
  IFS=';' read -ra dirs_to_backup <<< "$BACKUP_DIRS"
  for dir in ${dirs_to_backup[@]}; do

    dirname="$(echo "${dir}" | sed 's/\//-/g' | sed 's/^-//g')"
    mkdir "${host_dir}/${timestamp}_${dirname}"
    log_inner info "creating backup ${host_dir}/${timestamp}_${dirname} of ${dir}"

    # check for error codes and print error otherwise
    if ncftpget -R -v -u "$username" -p "$password" -P "${port}" "${address}" "${host_dir}/${timestamp}_${dirname}" "${dir}"; then

      # compress backup
      ( cd ${host_dir} && zip -r "${timestamp}_${dirname}.zip"  "${timestamp}_${dirname}") || log_inner error "error archiving ${host_dir}/${timestamp}_${dirname}"

      # removing downloaded files
      rm -fr "${host_dir}/${timestamp}_${dirname}"
      log_inner info "done backup ${host_dir}/${timestamp}_${dirname}"

    else
      log_inner error "error in downloading ${dir} from ${address}"
    fi

  done

  # setting status file to "backup done" to avoid consequent backups
  echo 0 > "$stat_file"
}

function reset(){

  if [[ -z $1 ]];then log_inner error "pass 3ds address as parameter"; return 1; else address="$1"; fi

  stat_file="${STAT_DIR}/${address}"
  # check if 3ds backup is running
  if [[ $(cat "$stat_file") == "2" ]]; then log_inner info "backup is running since $(stat -c '%y' "$STAT_FILE"), avoid resetting"; return 0; fi

  echo 1 > "$stat_file"
}

# main function that loops the given hosts and runs the backup script
function backup_cronjob(){

  IFS=';' read -ra addresses <<< "$FTPD_3DS_ADDRESSES"
  IFS=';' read -ra ports <<< "$FTPD_3DS_PORTS"
  IFS=';' read -ra usernames <<< "$FTPD_3DS_USERNAMES"
  IFS=';' read -ra passwords <<< "$FTPD_3DS_PASSWORDS"
  for index in "${!addresses[@]}"; do
    backup "${addresses[$index]}" "${ports[$index]}" "${usernames[$index]}" "${passwords[$index]}"
  done

}

function reset_cronjob(){

  IFS=';' read -ra addresses <<< "$FTPD_3DS_ADDRESSES"
  for index in "${!addresses[@]}"; do
    reset "${addresses[$index]}"
  done

}


case "$1" in
  "backup_cronjob")
    backup_cronjob
    ;;
  "reset_cronjob")
    reset_cronjob
    ;;
  *)
    log_inner error "usage: $0 backup_cronjob/reset_cronjob"
    ;;
esac
