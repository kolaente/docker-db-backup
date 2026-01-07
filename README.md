# Zero-Fuss Docker Database Backup

[![CI](https://github.com/kolaente/docker-db-backup/actions/workflows/ci.yml/badge.svg)](https://github.com/kolaente/docker-db-backup/actions/workflows/ci.yml)

A simple tool to create backup of all databases on a host. Supports postgres and mysql/mariadb.

Successor to [this script](https://kolaente.dev/konrad/docker-database-backup).

## Usage

Simply point it at your docker socket, mount a backup volume and be done:

```
docker run -v $PWD/backups:/backups -v /var/run/docker.sock:/var/run/docker.sock kolaente/db-backup
```

The tool will find all database containers running an official [`mysql`](https://hub.docker.com/_/mysql), 
[`mariadb`](https://hub.docker.com/_/mariadb) or [`postgres`](https://hub.docker.com/_/postgres) image and 
create backups of them periodically. It will also discover new containers 
as they are started and won't try to back up containers which have gone away.

When running, all backups for the current run are time-stamped into a sub folder of the backup directory (see below).

### Using labels

To make the backup tool discover other non-offical containers as well you can add the label `de.kolaente.db-backup` to 
any container with a value of `mysql` or `postgres` to treat it as a mysql or postgres container.

### Docker Compose

If you're running docker-compose, you can use a setup similar to the following compose file to run the backup:

```yaml
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

### `BACKUP_SCHEDULE`

The cron schedule at which the backup job runs, using the common unix cron syntax.

Check out [crontab.dev](https://crontab.dev/) for a nice explanation of the schedule.

Default: `0 */6 * * *` (every 6 hours)

### `BACKUP_NO_CRON`

If provided, runs the backup only once without a cron schedule. This is useful for one-off backups of testing if the configuration works.

Default: `false`

### `BACKUP_MAX`

How many backups to keep. If more backups are stored in the backup folder, the oldest one will be removed until there
are only as many as this config variable.

Default: `12`

### `BACKUP_COMPLETION_WEBHOOK_URL`

If provided, the tool will do an empty GET request to this URL to indicate it successfully completed the backup job.
You can use this with other tools to monitor if backups are completed as they should.

### `BACKUP_COMPRESS`

If set provided and set to `true`, all backups will be compressed using gzip.

## Building from source

This project uses go modules, so you'll need at least go 1.11 to compile it.

Simply run

```
go build .
```

to build the binary.
