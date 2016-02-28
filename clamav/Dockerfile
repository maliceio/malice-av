FROM gliderlabs/alpine

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-clamav
RUN apk-install clamav freshclam
RUN apk-install -t build-deps go git mercurial \
  && set -x \
  && mv /etc/clamav/freshclam.conf.sample /etc/clamav/freshclam.conf \
  && mv /etc/clamav/clamd.conf.sample /etc/clamav/clamd.conf \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-clamav \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scan \
  && rm -rf /go \
  && apk del --purge build-deps

ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Update ClamAV Definitions
RUN freshclam

WORKDIR /malware

ENTRYPOINT ["/bin/scan"]

CMD ["--help"]
