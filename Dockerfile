FROM golang:1-alpine AS build-env

# Setup repo
COPY . ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup
WORKDIR ${GOPATH}/src/kolaente.dev/konrad/docker-db-backup

RUN apk --no-cache add ca-certificates

RUN CGO_ENABLED=0 go build .

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /go/src/kolaente.dev/konrad/docker-db-backup /

CMD ["/docker-db-backup"]
