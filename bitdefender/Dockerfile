FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

# Install Requirements
RUN apt-get install -yq libc6-i386 wget && \
  wget -O- -q http://download.bitdefender.com/repos/deb/bd.key.asc | apt-key add - && \
  sh -c 'echo "deb http://download.bitdefender.com/repos/deb/ bitdefender non-free" >> /etc/apt/sources.list' && \
  apt-get update -qq

# Add Files
#ADD /run.sh /run.sh
#RUN chmod 755 /run.sh
# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Install Bitdefender
RUN DEBIAN_FRONTEND=noninteractive apt-get install -yq bitdefender-scanner
RUN cat /opt/BitDefender-scanner/etc/bdscan.conf > /opt/BitDefender-scanner/etc/bdscan.conf.bak
RUN echo "LicenseAccepted = True" >> /opt/BitDefender-scanner/etc/bdscan.conf
# Update Bitdefender
RUN bdscan --update

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

ENTRYPOINT ["bdscan"]

CMD ["--help"]
