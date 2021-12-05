FROM golang:1-alpine3.12 AS build-env

# Setup repo
COPY . ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup
WORKDIR ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup

RUN go build .

FROM alpine:3.12
LABEL maintainer="maintainers@kolaente.dev"

COPY --from=build-env /go/src/kolaente.dev/konrad/docker-db-backup /

CMD ["/docker-db-backup"]
