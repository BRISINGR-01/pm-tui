FROM python:3.12-slim

RUN apt-get update
RUN apt-get install -y wget git gcc sudo

# go env
ARG GO_VERSION="1.26.2"
ENV GO_VERSION=${GO_VERSION}

RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV PATH /go/bin:/usr/local/go/bin:$PATH

# project setup
WORKDIR /pm

COPY go.mod go.sum ./
# COPY tests/pip/package.json  ./

RUN go mod download

RUN pip install numpy panda

CMD go run .