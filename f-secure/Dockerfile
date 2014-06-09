FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

#Prevent daemon start during install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

# Install Requirements
RUN apt-get update -qq && apt-get install -yq libc6-i386

# Install F-Secure
ADD http://download.f-secure.com/webclub/f-secure-linux-security-10.00.60.tar.gz /

# Update F-Secure


# Add EICAR Test Virus File to malware folder
RUN mkdir malware && echo "X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*" > /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
