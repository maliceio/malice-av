FROM ubuntu:precise

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-comodo
RUN buildDeps='ca-certificates \
               build-essential \
               libssl-dev \
               golang \
               mercurial \
               git-core \
               gdebi \
               curl' \
  && apt-get install -yq python-software-properties \
  && add-apt-repository ppa:duh/golang \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
  && set -x \
  && echo "Install Comodo..." \
  && curl -Ls http://download.comodo.com/cavmgl/download/installs/1000/standalone/cav-linux_1.1.268025-1_amd64.deb > /tmp/cav-linux_1.1.268025-1_amd64.deb \
  && cd /tmp \
  && DEBIAN_FRONTEND=noninteractive gdebi -n cav-linux_1.1.268025-1_amd64.deb \
  && DEBIAN_FRONTEND=noninteractive /opt/COMODO/post_setup.sh \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-comodo \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/scan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps python-software-properties \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Update Comodo definitions
ADD http://download.comodo.com/av/updates58/sigs/bases/bases.cav /opt/COMODO/scanners/bases.cav
# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/scan"]

CMD ["--help"]
