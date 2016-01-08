FROM gliderlabs/alpine:3.3

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-clamav
RUN apk-install clamav
RUN apk-install -t build-deps go git mercurial \
	&& cd /go/src/github.com/maliceio/malice-gobased-plugin \
	&& export GOPATH=/go \
	&& go get \
	&& go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/clamscan \
	&& rm -rf /go \
	&& apk del --purge build-deps

ENTRYPOINT ["/bin/clamscan"]



# FROM debian:jessie
# MAINTAINER blacktop, https://github.com/blacktop
#
# RUN apt-get -qq update && apt-get install -yq clamav clamav-freshclam
#
# # Update ClamAV Definitions
# RUN freshclam
#
# # Add EICAR Test Virus File to malware folder
# ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR
#
# # Try to reduce size of container.
# RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
#
# WORKDIR /malware
#
# ENTRYPOINT ["/usr/bin/clamscan"]
#
# CMD ["--help"]
