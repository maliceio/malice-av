FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

# Install Requirements
RUN apt-get update -qq && apt-get install -yq libc6-i386

#Add F-PROT tar.gz
ADD http://files.f-prot.com/files/unix-trial/fp-Linux.x86.32-ws.tar.gz /
# RUN mv /fp-Linux.x86.32-ws.tar.gz/f-prot /opt/f-prot
RUN tar xvfz fp-Linux.x86.32-ws.tar.gz /opt/

# Install F-PROT
RUN ln -fs /opt/f-prot/fpscan /usr/local/bin/fpscan && \
  ln -fs /opt/f-prot/fpscand /usr/local/sbin/fpscand && \
  ln -fs /opt/f-prot/fpmon /usr/local/sbin/fpmon && \
  cp /opt/f-prot/f-prot.conf.default /opt/f-prot/f-prot.conf && \
  ln -fs /opt/f-prot/f-prot.conf /etc/f-prot.conf && \
  chmod a+x /opt/f-prot/fpscan && \
  chmod u+x /opt/f-prot/fpupdate && \
  ln -fs /opt/f-prot/man_pages/scan-mail.pl.8 /usr/share/man/man8/

# Update F-PROT Definitions
RUN /opt/f-prot/fpupdate

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

ENTRYPOINT ["/usr/local/bin/fpscan"]

CMD ["--help"]
