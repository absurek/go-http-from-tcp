
FROM ubuntu:24.04

ENV GOLANG_VERSION=1.25.5

# System setup
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    netcat-openbsd \
    git \
    jq \
    wget \
    ca-certificates \
    postgresql-client && \
    rm -rf /var/lib/apt/lists/*

# Go setup
RUN wget -q https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# GO CLIs
RUN go install github.com/bootdotdev/bootdev@latest && \
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
