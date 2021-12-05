FROM golang:1-alpine3.12 AS build-env

# Setup repo
COPY . ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup
WORKDIR ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup

RUN go build .

FROM alpine:3.12
LABEL maintainer="maintainers@kolaente.dev"

COPY --from=build-env /go/src/kolaente.dev/konrad/docker-db-backup /
COPY --from=postgres:14-alpine /usr/local/bin/pg_dump /usr/local/bin/pg_dump14
COPY --from=postgres:13-alpine /usr/local/bin/pg_dump /usr/local/bin/pg_dump13
COPY --from=postgres:12-alpine /usr/local/bin/pg_dump /usr/local/bin/pg_dump12
COPY --from=postgres:11-alpine /usr/local/bin/pg_dump /usr/local/bin/pg_dump11
COPY --from=postgres:10-alpine /usr/local/bin/pg_dump /usr/local/bin/pg_dump10
COPY --from=mariadb:10 /usr/bin/mysqldump /usr/local/bin/mysqldump

CMD ["/docker-db-backup"]
