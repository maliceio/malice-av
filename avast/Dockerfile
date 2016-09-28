# FROM ubuntu:trusty
FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

# Install Avast AV
COPY license.avastlic /etc/avast/license.avastlic
ENV AVAST_VERSION 2.0.0-1
# ENV AVAST_VERSION 2.1.1-1

RUN \
    echo "Install Avast..." \
    && echo 'deb http://deb.avast.com/lin/repo debian release' >> /etc/apt/sources.list \
    && apt-key adv --fetch-keys https://files.avast.com/files/resellers/linux/avast.gpg \
    && apt-get update -q && apt-get install -y libc6-i386 avast=$AVAST_VERSION \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Go binary
COPY . /go/src/github.com/maliceio/malice-avast
RUN buildDeps='build-essential \
               mercurial \
               git-core \
               wget' \
    && apt-get update -qq \
    && apt-get install -yq $buildDeps --no-install-recommends \
    && echo "Install Go..." \
    && ARCH="$(dpkg --print-architecture)" \
    && wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && export PATH=$PATH:/usr/local/go/bin \
    && echo "Building avscan Go binary..." \
    && cd /go/src/github.com/maliceio/malice-avast \
    && export GOPATH=/go \
    && go version \
    && go get \
    && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan \
    && echo "Clean up unnecessary files..." \
    && apt-get purge -y --auto-remove $buildDeps \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go /usr/local/go

COPY eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]

CMD ["--help"]

# NOTE: https://www.avast.com/en-us/faq.php?article=AVKB131
# NOTE: To Update run - /var/lib/avast/Setup/avast.vpsupdate
