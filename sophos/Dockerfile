FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

# Install Requirements
RUN buildDeps='ca-certificates wget' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
  && echo "Install Sophos..." \
  && cd /tmp \
  && wget https://github.com/maliceio/malice-av/raw/master/sophos/sav-linux-free-9.tgz \
  && tar xzvf sav-linux-free-9.tgz \
  && ./sophos-av/install.sh /opt/sophos --update-free --acceptlicence --autostart=False --enableOnBoot=False --automatic --ignore-existing-installation --update-source-type=s \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Install Go binary
COPY . /go/src/github.com/maliceio/malice-sophos
RUN buildDeps='ca-certificates \
               build-essential \
               golang-go \
               mercurial \
               git-core' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps --no-install-recommends \
  && echo "Building sophos malscan Go binary..." \
  && cd /go/src/github.com/maliceio/malice-sophos \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/malscan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Update Sophos
RUN /opt/sophos/update/savupdate.sh

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/malscan"]

CMD ["--help"]
