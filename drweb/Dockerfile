FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

RUN buildDeps='ca-certificates \
               build-essential \
               libfontconfig1 \
               libxrender1 \
               libglib2.0-0 \
               libxi6 \
               xauth \
               gdebi-core \
               libssl-dev \
               mercurial \
               git-core \
               wget' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 \
  && set -x \
  && echo "Install Dr Web..." \
  && echo 'deb http://repo.drweb.com/drweb/debian 10.0.0 non-free' >> /etc/apt/sources.list \
  && apt-key adv --fetch-keys http://repo.drweb.com/drweb/drweb.key \
  && apt-get update -q && apt-get install -y drweb-workstations \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Go binary
COPY . /go/src/github.com/maliceio/malice-drweb
RUN buildDeps='build-essential \
               golang-go \
               mercurial \
               git-core' \
    && apt-get update -qq \
    && apt-get install -yq $buildDeps --no-install-recommends \
    && echo "Building drweb malscan Go binary..." \
    && cd /go/src/github.com/maliceio/malice-drweb \
    && export GOPATH=/go \
    && go version \
    && go get \
    && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/malscan \
    && echo "Clean up unnecessary files..." \
    && apt-get purge -y --auto-remove $buildDeps \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Update Dr Web AV Definitions
# RUN drweb-ctl update

WORKDIR /malware

ENTRYPOINT ["/bin/malscan"]

CMD ["--help"]
