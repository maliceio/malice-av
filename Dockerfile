FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get install -yq wget

ADD http://www.microworldsystems.com/download/linux/5.5/eScan/ubuntu12.04/64b/escan-5.5-2.Ubuntu.12.04_x86_64.deb
