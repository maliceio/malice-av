FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

RUN buildDeps='ca-certificates \
               wget' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps iptables-persistent libc6-i386 \
  && set -x \
  && echo "Install Dr Web..." \
  && cd /tmp \
  && wget http://download.geo.drweb.com/pub/drweb/unix/workstation/11.0/drweb-11.0.1-av-linux-amd64.run \
  && chmod +x drweb-11.0.1-av-linux-amd64.run
  # && echo 'deb http://repo.drweb.com/drweb/debian 10.0.0 non-free' >> /etc/apt/sources.list \
  # && apt-key adv --fetch-keys http://repo.drweb.com/drweb/drweb.key \
  # && apt-get update -q && apt-get install -y drweb-workstations \
  # && apt-get purge -y --auto-remove $buildDeps \
  # && apt-get clean \
  # && rm -rf /var/lib/apt/lists/* /tmp/*

#   && cd /tmp \
#   && wget http://download.geo.drweb.com/pub/drweb/unix/workstation/11.0/drweb-11.0.1-av-linux-amd64.run \
#   && chmod +x drweb-11.0.1-av-linux-amd64.run
# http://download.geo.drweb.com/pub/drweb/unix/workstation/11.0/drweb-11.0.1-av-linux-amd64.run

#                build-essential \
#                libfontconfig1 \
#                libxrender1 \
#                libglib2.0-0 \
#                libxi6 \
#                xauth \
#                gdebi-core \
#                libssl-dev \
#                mercurial \
#                git-core \
#
# libc6-i386

# Install Go binary
# COPY . /go/src/github.com/maliceio/malice-drweb
# RUN buildDeps='ca-certificates \
#                build-essential \
#                golang-go \
#                mercurial \
#                git-core' \
#   && apt-get update -qq \
#   && apt-get install -yq $buildDeps --no-install-recommends \
#   && echo "Building drweb malscan Go binary..." \
#   && cd /go/src/github.com/maliceio/malice-drweb \
#   && export GOPATH=/go \
#   && go version \
#   && go get \
#   && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/malscan \
#   && echo "Clean up unnecessary files..." \
#   && apt-get purge -y --auto-remove $buildDeps \
#   && apt-get clean \
#   && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Update Dr Web AV Definitions
# RUN drweb-ctl update

# WORKDIR /malware
#
# ENTRYPOINT ["/bin/malscan"]
#
# CMD ["--help"]
