FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

ENV GO_VERSION 1.7.1
ENV FSECURE_VERSION 11.00.79-rtm

# Install Requirements
RUN buildDeps='ca-certificates wget rpm' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps lib32stdc++6 psmisc \
  && echo "Install F-Secure..." \
  && cd /tmp \
  && wget https://download.f-secure.com/corpro/ls/trial/fsls-${FSECURE_VERSION}.tar.gz \
  && tar zxvf fsls-${FSECURE_VERSION}.tar.gz \
  && cd fsls-${FSECURE_VERSION} \
  && chmod a+x fsls-${FSECURE_VERSION} \
  && ./fsls-${FSECURE_VERSION} --auto standalone lang=en --command-line-only \
  && fsav --version \
  && echo "Update F-Secure..." \
  && cd /tmp \
  && wget http://download.f-secure.com/latest/fsdbupdate9.run \
  && mv fsdbupdate9.run /opt/f-secure/ \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/*

# Update F-Secure
RUN /etc/init.d/fsaua start \
  && /etc/init.d/fsupdate start \
  && /opt/f-secure/fsav/bin/dbupdate /opt/f-secure/fsdbupdate9.run; exit 0

# Install Go binary
COPY . /go/src/github.com/maliceio/malice-fsecure
RUN buildDeps='ca-certificates \
               build-essential \
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
  && cd /go/src/github.com/maliceio/malice-fsecure \
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

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]

CMD ["--help"]
