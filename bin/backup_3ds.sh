#!/bin/bash
KEEP_LAST=${KEEP_LAST:-10}
BASE_DIR=${BASE_DIR:-/var/lib/backup_3ds}
BACKUP_DEST=${BACKUP_DEST:-$BASE_DIR/backups}
TEMPLATE_DIR=${TEMPLATE_DIR:-/var/lib/backup_3ds/dashboard/templates}
WEB_ROOT=${WEB_ROOT:-/var/lib/backup_3ds/dashboard/static}
CONFIG_FILE=${CONFIG_FILE:-/etc/backup_3ds/config.yml}
STAT_DIR=${STAT_DIR:-/tmp/backup_3ds/status}

# Check if the config file exists
test -f "$CONFIG_FILE" || { echo "Config file $CONFIG_FILE does not exist. Please create it."; exit 1; }


mkdir -p "$BASE_DIR" "$BACKUP_DEST" "$TEMPLATE_DIR" "$WEB_ROOT" "$STAT_DIR"

log_inner() {
  echo "[${BASH_SOURCE##*/}:${1^^}] ${FUNCNAME[2]}@${BASH_LINENO[1]}: ${*:2}"
}

function telegram_send() {
  test -z $TELEGRAM_TOKEN && log_inner warning "TELEGRAM_TOKEN is not set" && return 1
  test -z $TELEGRAM_CHAT_ID && log_inner warning "TELEGRAM_CHAT_ID is not set" && return 1

  MESSAGE="$(cat)"

  RESPONSE="$(curl -sSf \
    --request POST \
    --url "https://api.telegram.org/bot${TELEGRAM_TOKEN}/sendMessage" \
    --header "Content-Type: application/json" \
    --data "$(jq -n \
      --arg chat_id "$TELEGRAM_CHAT_ID" \
      --arg text "$MESSAGE" \
      '{chat_id: $chat_id, text: $text}'
    )")"

if echo "$RESPONSE" | jq -e '.ok == true' > /dev/null; then
  log_inner info "Message to $TELEGRAM_CHAT_ID sent successfully."
else
  log_inner warning "Error sending message to $TELEGRAM_CHAT_ID" >&2
  echo "$RESPONSE" | jq '.description' >&2
  return 1
fi
}

# backup a single 3ds, parameters address port username password
function backup(){


  if [[ -z $1 ]] || [[ "$1" == "null" ]];then log_inner error "set console name"; return 4; else name="$1"; fi
  if [[ -n $2 ]] && [[ "$2" != "null" ]];then port="$2"; else port=21; fi
  if [[ -n $3 ]] && [[ "$3" != "null" ]];then user="$3"; fi
  if [[ -n $4 ]] && [[ "$4" != "null" ]];then password="$4"; fi

  if [[ -n $user ]] && [[ -n $password ]]; then userpass="-u ${user},${password}"; else userpass=""; fi

  # set status file variable and initialize the file if it does not exist
  stat_file="${STAT_DIR}/${name}"

  if [[ ! -f "$stat_file" ]]; then echo 1 > "$stat_file"; fi

  # setting backup dir for the specific 3ds
  host_dir="${BACKUP_DEST}/${name}"; mkdir -p "$host_dir"

  # check if 3ds backup is running, exit if yes
  if [[ $(cat "$stat_file") == "2" ]]; then log_inner info "backup of ${name} is running since $(stat -c '%y' "$stat_file")"; return 2; fi

  # check if 3ds backup has been made, exit if yes
  if [[ $(cat "$stat_file") == "0" ]]; then log_inner info "backup of ${name} has already been made at $(stat -c '%y' "$stat_file")"; return 1; fi

  # check if 3ds ftp server is up, exit if not
  if ! nc -z -w1 "$name" "$port"; then log_inner info "console at ${name}:${port} is not listening for ftp connections"; return 3; fi

  # setting lock file to avoid running multiple backup jobs in parallel on the same 3DS
  echo 2 > "$stat_file"
  log_inner info "starting console backup of $name:$port"
  echo -e "Starting ${name}'s backup of\n$(yq ".consoles[] | select(.name == \"$name\").dirs | join(\"\n\")" "${CONFIG_FILE}" -r)" | telegram_send

  timestamp="$(date +%s)"

  # loop all dirs to backup
  yq ".consoles[] | select(.name == \"$name\").dirs | join(\"\n\")" "${CONFIG_FILE}" -r | while read dir; do

    dirname="$(echo "${dir}" | sed 's/\//-/g' | sed 's/^-//g')"
    mkdir "${host_dir}/${name}_${timestamp}_${dirname}"
    log_inner info "creating backup ${host_dir}/${name}_${timestamp}_${dirname} of ${dir}"

    # check for error codes and print error otherwise

    log_inner info lftp -e "mirror --verbose=3 ${dir} ${host_dir}/${name}_${timestamp}_${dirname}" -p "${port}" "${userpass}" "${name}"
    lftp -e "mirror --verbose=3 ${dir} ${host_dir}/${name}_${timestamp}_${dirname}" -p "${port}" "${userpass}" "${name}"

    # compress backup
    ( cd ${host_dir} && zip -r "${name}_${timestamp}_${dirname}.zip"  "${name}_${timestamp}_${dirname}") || log_inner error "error archiving ${host_dir}/${name}_${timestamp}_${dirname}"

    # removing downloaded files
    rm -fr "${host_dir}/${name}_${timestamp}_${dirname}"
    log_inner info "done backup of ${host_dir}/${name}_${timestamp}_${dirname} from ${name}"

  done

  # setting status file to "backup done" to avoid consequent backups
  echo -e "Done ${name}'s backup of\n$(yq ".consoles[] | select(.name == \"$name\").dirs | join(\"\n\")" "${CONFIG_FILE}" -r)" | telegram_send
  echo 0 > "$stat_file"
}

function reset(){

  if [[ -z $1 ]] || [[ "$1" == "null" ]];then log_inner error "set console name"; return 4; else name="$1"; fi

  stat_file="${STAT_DIR}/${name}"
  # check if 3ds backup is running
  if [[ $(cat "$stat_file") == "2" ]]; then log_inner info "backup for $address is running since $(stat -c '%y' "$STAT_FILE"), avoid resetting"; return 0; fi

  echo 1 > "$stat_file"
}

function generate_dashboard(){

  # generate index
  source <( bash-tpl "$TEMPLATE_DIR/index.html.tpl" ) > "$WEB_ROOT/index.html"

  # generate backups pages for each console
  for dir in $(find "$BACKUP_DEST/"  -maxdepth 1 -not -path "$BACKUP_DEST/" -type d); do
    dir_name="$(basename "$dir")"
    log_inner info "generating backup page for $dir_name using data from $dir"
    log_inner info "source <( bash-tpl "$TEMPLATE_DIR/backup_list.html.tpl" ) > $WEB_ROOT/$dir_name.html"
    source <( bash-tpl "$TEMPLATE_DIR/backup_list.html.tpl" ) > "$WEB_ROOT/$dir_name.html"
  done

}


prune() {

  # test if the console name is provided
  if [[ -z $1 ]] || [[ "$1" == "null" ]];then log_inner error "set console name"; return 4; else name="$1"; fi

  # get keep_last value from config file for the specific console, if not set use default KEEP_LAST
  keep_last=$(yq ".consoles[] | select(.name == \"$name\").keep_last" "${CONFIG_FILE}")
if [[ -z $keep_last ]] || [[ "$keep_last" == "null" ]]; then log_inner info "KEEP_LAST not set for $name, using default value of $KEEP_LAST"; keep_last=$KEEP_LAST; fi

  # Validate KEEP_LAST is a positive integer
  if ! [[ "$keep_last" =~ ^[0-9][0-9]*$ ]]; then log_inner error "KEEP_LAST must be a positive integer, but got $keep_last"; exit 1; fi

  # avoid pruning backups if KEEP_LAST variable is 0
  if [[ $keep_last == "0" ]];then log_inner info " KEEP_LAST set to 0, avoid pruning "; return 0; fi

  log_inner info "Running pruning job for $name, keeping the last $KEEP_LAST."

  # Find all immediate subdirectories within BACKUP_DEST,
  # sort them by modification time (oldest first).
  # Then calculate how many to delete to keep only the KEEP_LAST newest.
  backups_list=$(find "$BACKUP_DEST/$name" -name '*.zip' -printf '%T@ %p\n' | sort -n)
  total_backups=$(echo "$backups_list" | wc -l)

  # Calculate how many backups to prune
  num_to_prune=$(( total_backups - keep_last ))

  if [ "$num_to_prune" -gt 0 ]; then
    log_inner info "Identified $total_backups backups in total, keeping $keep_last. Pruning $num_to_prune oldest backups."
    echo "$backups_list" | \
      head -n "$num_to_prune" | \
      cut -d' ' -f2- | \
      xargs -r rm -rf

    if [ $? -eq 0 ]; then
      log_inner info "Pruning process completed successfully."
    else
      log_inner error "Pruning process encountered an error while removing files."
      return 1
    fi
  else
    log_inner info "No backups to prune. Total backups: $total_backups, desired to keep: $keep_last."
  fi
}

# main function that loops the given hosts and runs the backup script
function backup_cronjob(){

  log_inner info "Starting backup cronjob for all consoles defined in $CONFIG_FILE"
  yq '.consoles[] | "\(.name) \(.port) \(.user) \(.password)"' "${CONFIG_FILE}" -r | while read name port user password; do
    log_inner info "Parsed console from config: $name $port $user $password"
    backup "$name" "$port" "$user" "$password" && generate_dashboard
  done

}

function reset_cronjob(){

  yq '.consoles[] | "\(.name)"' "${CONFIG_FILE}" -r | while read name; do
    reset "${name}"
  done

}

function prune_cronjob(){

  log_inner info "Starting prune cronjob for all consoles defined in $CONFIG_FILE"
  yq '.consoles[] | "\(.name)"' "${CONFIG_FILE}" -r | while read name; do
    prune "${name}"
  done

}

case "$1" in
  "backup_cronjob")
    backup_cronjob
    ;;
  "reset_cronjob")
    reset_cronjob
    ;;
  "prune_cronjob")
    prune_cronjob
    ;;
  *)
    log_inner error "usage: $0 backup_cronjob/reset_cronjob/prune_cronjob"
    ;;
esac
