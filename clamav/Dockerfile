FROM gliderlabs/alpine:3.3

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-clamav
RUN apk-install clamav clamav-dev
RUN apk-install -t build-deps go git mercurial \
	&& cd /go/src/github.com/maliceio/malice-clamav \
	&& export GOPATH=/go \
	&& go get \
	&& go build -ldflags "-X main.Version=$(cat VERSION)" -o /bin/clamscan \
	&& rm -rf /go \
	&& apk del --purge build-deps

RUN mv /etc/clamav/freshclam.conf.sample /etc/clamav/freshclam.conf \
  && mv /etc/clamav/clamd.conf.sample /etc/clamav/clamd.conf \
	&& freshclam

ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

ENTRYPOINT ["/bin/clamscan"]
