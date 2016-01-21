FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-avira
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
  && echo "Install Avira..." \
  && curl -Ls http://premium.avira-update.com/package/wks_avira/unix/en/pers/antivir_workstation-pers.tar.gz > /tmp/avira.tar.gz \
  && cd /tmp \
  && tar -zxvf avira.tar.gz \
  && antivir-workstation-pers-3.1.3.5-0/install --inf=/go/src/github.com/maliceio/malice-avira/unattended.inf \
  && curl -Ls http://personal.avira-update.com/package/peclkey/win32/int/hbedv.key > /usr/lib/AntiVir/guard/avira.key \
  && mkdir /home/quarantine/ \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-avira \
  && export GOPATH=/go \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Update Avira
RUN /usr/lib/AntiVir/guard/avupdate-guard --product=Guard

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

# ENTRYPOINT ["/usr/local/bin/fpscan"]
ENTRYPOINT ["/bin/scan"]

CMD ["--help"]

# FROM ubuntu:latest
# MAINTAINER blacktop, https://github.com/blacktop
#
# RUN apt-get -q update && apt-get install -yq libc6-i386
#
# # Add Files
# ADD /run.sh /run.sh
# RUN chmod 755 /run.sh
# # Add EICAR Test Virus File to malware folder
# ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR
# ADD /unattended.inf /unattended.inf
# #ADD /supervisord.conf /etc/supervisor/conf.d/supervisord.conf
# RUN mkdir /home/quarantine/ #&& mkdir -p /var/log/supervisor
#
# # Download Avira
# ADD http://premium.avira-update.com/package/wks_avira/unix/en/pers/antivir_workstation-pers.tar.gz /
# RUN tar -zxvf /antivir_workstation-pers.tar.gz
#
# # Install Avira
# RUN /antivir-workstation-pers-3.1.3.5-0/install --inf=/unattended.inf
# ADD http://personal.avira-update.com/package/peclkey/win32/int/hbedv.key /usr/lib/AntiVir/guard/avira.key
# # Update Avira
# RUN /usr/lib/AntiVir/guard/avupdate-guard --product=Guard
#
# # Try to reduce size of container.
# RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
#
# WORKDIR /malware
#
# CMD ["/run.sh"]
