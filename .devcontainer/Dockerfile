ARG GO_VERSION=1.22.1

FROM golang:${GO_VERSION}

RUN apt-get  update -y -q && \
    apt-get upgrade -y -q

RUN apt-get install -y sqlite3 libsqlite3-dev
