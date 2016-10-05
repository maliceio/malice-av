FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

ENV GO_VERSION 1.7.1

COPY . /go/src/github.com/maliceio/malice-fprot
RUN buildDeps='ca-certificates \
               build-essential \
               mercurial \
               git-core \
               unzip \
               wget' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 --no-install-recommends \
  && set -x \
  && echo "Install F-PROT..." \
  && wget https://github.com/maliceio/malice-av/raw/master/fprot/fp-Linux.x86.32-ws.tar.gz \
    -O /tmp/fp-Linux.x86.32-ws.tar.gz \
  && tar -C /opt -zxvf /tmp/fp-Linux.x86.32-ws.tar.gz \
  && ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan \
  && ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand \
  && ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon \
  && cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf \
  && ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf \
  && chmod a+x /opt/f-prot/fpscan \
  && chmod u+x /opt/f-prot/fpupdate \
  && ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/ \
  && echo "Install Go..." \
  && ARCH="$(dpkg --print-architecture)" \
  && wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz \
  && tar -C /usr/local -xzf /tmp/go.tar.gz \
  && export PATH=$PATH:/usr/local/go/bin \
  && echo "Building avscan Go binary..." \
  && cd /go/src/github.com/maliceio/malice-fprot \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go /usr/local/go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Update F-PROT Definitions
RUN /opt/f-prot/fpupdate

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]

CMD ["--help"]

# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.64-fs.tar.gz
# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.32-ws.tar.gz
# http://files.f-prot.com/files/unix-trial/fp-Linux.x86.64-ws.tar.gz
