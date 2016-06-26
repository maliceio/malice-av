FROM gliderlabs/alpine

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-clamav
RUN apk-install clamav freshclam ca-certificates
RUN apk-install -t build-deps go git mercurial \
  && set -x \
  && echo "Building scan Go binary..." \
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
