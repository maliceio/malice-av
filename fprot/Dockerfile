FROM debian:wheezy
MAINTAINER blacktop, https://github.com/blacktop

RUN buildDeps='ca-certificates \
               build-essential \
               golang-go \
               mercurial \
               git-core \
               unzip \
               curl' \
  && set -x \
  && echo 'deb http://ftp.us.debian.org/debian/ testing main contrib non-free' >> /etc/apt/sources.list \
  && apt-get update -qq \
  && apt-get install -t testing libc6 \
  && apt-get install -yq $buildDeps libimage-exiftool-perl --no-install-recommends \
  && set -x \
  && echo "Install F-PROT..." \
  && curl -Ls http://files.f-prot.com/files/unix-trial/fp-Linux.x86.32-ws.tar.gz > /tmp/fprot.tar.gz \
  && cd /tmp \
  && tar -C /opt -zxvf fprot.tar.gz \
  && ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan \
  && ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand \
  && ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon \
  && cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf \
  && ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf \
  && chmod a+x /opt/f-prot/fpscan \
  && chmod u+x /opt/f-prot/fpupdate \
  && ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/ \
  && echo "Updating F-PROT Definitions..." \
  && /opt/f-prot/fpupdate \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-fileinfo \
  && export GOPATH=/go \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/info \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/usr/local/bin/fpscan"]

CMD ["--help"]
