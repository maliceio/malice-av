FROM debian:jessie

MAINTAINER blacktop, https://github.com/blacktop

ENV FSECURE_VERSION 11.00.79-rtm

# Install Requirements
RUN buildDeps='ca-certificates lib32stdc++6 wget rpm' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
  && echo "Install F-Secure..." \
  && cd /tmp \
  && wget https://download.f-secure.com/corpro/ls/trial/fsls-${FSECURE_VERSION}.tar.gz \
  && tar zxvf fsls-${FSECURE_VERSION}.tar.gz \
  && cd fsls-${FSECURE_VERSION} \
  && chmod a+x fsls-${FSECURE_VERSION} \
  && ./fsls-${FSECURE_VERSION} --auto lang=en --command-line-only \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# RUN printf "f\nw\n\n" | /opt/f-secure/fsav/sbin/fschooser

# Install Go binary
COPY . /go/src/github.com/maliceio/malice-fsecure
RUN buildDeps='ca-certificates \
               build-essential \
               golang-go \
               mercurial \
               git-core' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps --no-install-recommends \
  && echo "Building drweb malscan Go binary..." \
  && cd /go/src/github.com/maliceio/malice-fsecure \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version $(cat VERSION) -X main.BuildTime $(date -u +%Y%m%d)" -o /bin/malscan \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go

# Update F-Secure
RUN set -x \
  && /etc/init.d/fsaua start \
  && /etc/init.d/fsupdated start \
  &&
# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/malscan"]

CMD ["--help"]
