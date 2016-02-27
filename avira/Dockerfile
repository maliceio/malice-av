FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-avira
RUN buildDeps='ca-certificates \
               build-essential \
               mercurial \
               git-core \
               golang \
               unzip \
               curl' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps libc6-i386 file --no-install-recommends \
  && set -x \
  && echo "Install Avira..." \
  && cd /tmp \
  && curl -O http://premium.avira-update.com/package/wks_avira/unix/en/pers/antivir_workstation-pers.tar.gz \
  && tar -zxvf antivir_workstation-pers.tar.gz \
  && antivir-workstation-pers-3.1.3.5-0/install --inf=/go/src/github.com/maliceio/malice-avira/unattended.inf \
  && mkdir /home/quarantine/ \
  && echo "Building malice-avira Go binary..." \
  && cd /go/src/github.com/maliceio/malice-avira \
  && export GOPATH=/go \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/scan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Update Avira
COPY hbedv.key /usr/lib/AntiVir/guard/avira.key
RUN /usr/lib/AntiVir/guard/avupdate-guard --product=Guard

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/scan"]

CMD ["--help"]
