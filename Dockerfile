FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

# Make sure that Upstart won't try to start avgd after dpkg installs it
# https://github.com/dotcloud/docker/issues/446
ADD policy-rc.d /usr/sbin/policy-rc.d
RUN /bin/chmod 755 /usr/sbin/policy-rc.d

# Install Requirements
RUN apt-get update -qq
RUN apt-get install -qqy libc6-i386 supervisor wget

# Get AVG Installer
#ADD http://files.avast.com/files/linux/avast4workstation_1.3.0-2_i386.deb /avast/
ADD avast4workstation_1.3.0-2_i386.deb /avast/avast4workstation_1.3.0-2_i386.deb

# Install AVG
RUN dpkg -i /avast/avast4workstation_1.3.0-2_i386.deb

# Start Avast Service and Update Avast Definitions
ADD avastrc .avast/avastrc
WORKDIR /malware
RUN sysctl -w kernel.shmmax=128000000
RUN /usr/bin/avast && /usr/bin/avast-update

# Add Files
ADD /supervisord-avast.conf /etc/supervisor/conf.d/supervisord-avast.conf
ADD /malware/EICAR /malware/EICAR
ADD /run.sh /run.sh
RUN chmod 755 /*.sh

ENTRYPOINT ["/usr/bin/avast"]
CMD ["--help"]
