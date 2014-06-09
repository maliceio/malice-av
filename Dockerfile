FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

#Prevent daemon start during install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

# Install Requirements
RUN apt-get update -qq && apt-get install -qqy libc6-i386 supervisor wget

# Add EICAR Test Virus File to malware folder
RUN mkdir malware && echo "X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*" > /malware/EICAR

# Get Avast Installer
ADD http://files.avast.com/files/linux/avast4workstation_1.3.0-2_i386.deb /avast/

# Install Avast
RUN dpkg -i /avast/avast4workstation_1.3.0-2_i386.deb

# Start Avast Service and Update Avast Definitions
ADD avastrc .avast/avastrc

RUN sysctl -w kernel.shmmax=128000000
RUN /usr/bin/avast && /usr/bin/avast-update

WORKDIR /malware

ENTRYPOINT ["/usr/bin/avast"]

CMD ["--help"]
