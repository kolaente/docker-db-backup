FROM golang:1-alpine3.12 AS build-env

# Setup repo
COPY . ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup
WORKDIR ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup

RUN CGO_ENABLED=0 go build .

FROM scratch

COPY --from=build-env /go/src/kolaente.dev/konrad/docker-db-backup /

CMD ["/docker-db-backup"]
