FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

# Install Requirements
COPY . /go/src/github.com/maliceio/malice-avg
RUN buildDeps='ca-certificates \
               build-essential \
               golang-go \
               mercurial \
               git-core \
               unzip \
               curl' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && set -x \
  && echo "Install AVG..." \
  && curl -Ls http://download.avgfree.com/filedir/inst/avg2013flx-r3118-a6926.i386.deb > /tmp/avg.deb \
  && dpkg -i /tmp/avg.deb \
  && /etc/init.d/avgd restart && avgupdate \
  && avgcfgctl -w UpdateVir.sched.Task.Disabled=true \
  && avgcfgctl -w Default.setup.daemonize=false   \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-avg \
  && export GOPATH=/go \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/scan"]

CMD ["--help"]
