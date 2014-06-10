FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

#Prevent daemon start during install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

# Install Requirements
RUN apt-get update && apt-get install -y gdebi
# Get Comodo Installer
ADD http://download.comodo.com/cavmgl/download/installs/1000/standalone/cav-linux_1.1.268025-1_amd64.deb /comodo/
# Install Comodo
RUN DEBIAN_FRONTEND=noninteractive gdebi -n /comodo/cav-linux_1.1.268025-1_amd64.deb
# Run Post Install Script and Accept Defaults
RUN DEBIAN_FRONTEND=noninteractive /opt/COMODO/post_setup.sh
# Update Comodo AV Definitions
ADD http://download.comodo.com/av/updates58/sigs/bases/bases.cav /opt/COMODO/scanners/bases.cav

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

ENTRYPOINT ["/opt/COMODO/cmdscan"]

CMD ["-h"]
