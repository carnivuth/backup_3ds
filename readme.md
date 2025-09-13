# Backup 3ds

Make backups of 3ds using ftpd periodically

## Installation

Installation is done trough docker container a sample can be as follows

```yaml

services:
  backup_3ds:
    container_name: backup_3ds
    image: carnivuth/backup_3ds:latest
    environment:
      - FTPD_3DS_ADDRESS=${FTPD_3DS_ADDRESS}
      - FTPD_3DS_PORT=${FTPD_3DS_PORT}
    volumes:
      - "./data:/var/lib/3ds_backup"
```
