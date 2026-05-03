FROM archlinux:latest

# Update and install base dependencies
RUN pacman -Syu --noconfirm \
    wget base-devel \
    git \
    sudo
    
# go env
ARG GO_VERSION="1.26.2"
ENV GO_VERSION=${GO_VERSION}

RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV PATH /go/bin:/usr/local/go/bin:$PATH


# Create a non-root user
RUN useradd -m -G wheel -s /bin/bash builder \
    && echo "builder ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER builder
WORKDIR /home/builder

# Install yay
RUN git clone https://aur.archlinux.org/yay.git \
    && cd yay \
    && makepkg -si --noconfirm

# project setup
WORKDIR /pm

COPY go.mod go.sum ./

RUN go mod download

CMD go run .