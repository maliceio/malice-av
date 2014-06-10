FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get -qq update && apt-get install -yq clamav clamav-freshclam

# Update ClamAV Definitions
RUN freshclam

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

ENTRYPOINT ["/usr/bin/clamscan"]

CMD ["--help"]
