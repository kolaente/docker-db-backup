# Zero-Fuss Docker Database Backup

[![Build Status](https://drone.kolaente.de/api/badges/konrad/docker-db-backup/status.svg?ref=refs/heads/main)](https://drone.kolaente.de/konrad/docker-db-backup)

A simple tool to create backup of all databases on a host. Supports postgres and mysql/mariadb.

Successor to [this script](https://kolaente.dev/konrad/docker-database-backup).

## Usage

Simply point it at your docker socket, mount a backup volume and be done:

```
docker run -v $PWD:/backups -v /var/run/docker.sock:/var/run/docker.sock kolaente/db-backup
```

The tool will find all database containers and create backups of them periodically. It will also discover new containers
as they are started and won't try to back up containers which have gone away.

When running, all backups for the current run are time-stamped into a sub folder of the backup directory (see below).

### Docker Compose

If you're running docker-compose, you can use a setup similar to the following compose file to run the backup:

```yaml
version: '2'
services:
  backup:
    image: kolaente/db-backup
    restart: unless-stopped
    volumes:
      - ./backups:/backups
      - /etc/localtime:/etc/localtime:ro
      - /var/run/docker.sock:/var/run/docker.sock

```

## Config

All config is done with environment variables.

### `BACKUP_FOLDER`

Where all backup files will be stored.

Default: `/backups`

### `BACKUP_INTERVAL`

The interval at which backups will happen. Must be a parsable string
as [time.Duration](https://pkg.go.dev/time#ParseDuration). Must be positive.

Default: `3h`

### `BACKUP_MAX`

How many backups to keep. If more backups are stored in the backup folder, the oldest one will be removed until there
are only as many as this config variable.

Default: `12`

## Building from source

This project uses go modules, so you'll need at least go 1.11 to compile it.

Simply run

```
go build .
```

to build the binary.
