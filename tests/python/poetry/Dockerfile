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

# poetry setup
ENV POETRY_VERSION=1.8.3 \
    POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALENVS_CREATE=false

RUN pip install --no-cache-dir poetry

# project setup
WORKDIR /pm


COPY go.mod go.sum ./
RUN go mod download

COPY tests/python/poetry/pyproject.toml tests/python/poetry/pyproject.toml  
WORKDIR tests/python/poetry/
RUN poetry install --no-root 

# CMD go run ..
CMD bash