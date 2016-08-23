FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

#Prevent daemon start during install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

RUN apt-get install -yq libc6-i386 gdebi wget

ADD http://www.microworldsystems.com/download/linux/5.5/eScan/ubuntu12.04/64b/escan-5.5-2.Ubuntu.12.04_x86_64.deb

# Install EScan
RUN DEBIAN_FRONTEND=noninteractive gdebi -n escan-5.5-2.Ubuntu.12.04_x86_64.deb

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# http://www.microworldsystems.com/download/linux/soho/deb/escan-antivirus-wks-7.0-18.amd64.deb


# http://www.microworldsystems.com/download/linux/soho/deb/escan-antivirus-wks-7.0-20.amd64.deb
