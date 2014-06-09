FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get -qq update && apt-get install -yq clamav clamav-freshclam

# Update ClamAV Definitions
RUN freshclam

# Add EICAR Test Virus File to malware folder
RUN mkdir malware && echo "X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*" > /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

ENTRYPOINT ["/usr/bin/clamscan"]

CMD ["--help"]
