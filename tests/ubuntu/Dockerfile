FROM ubuntu:latest

ARG GO_VERSION="1.26.2"
ENV GO_VERSION=${GO_VERSION}

RUN apt-get update
RUN apt-get install -y wget git gcc sudo

RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV PATH /go/bin:/usr/local/go/bin:$PATH

WORKDIR /pm

COPY go.mod go.sum ./

RUN go mod download

CMD go run .