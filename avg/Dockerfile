FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

#Prevent daemon start during install
RUN echo '#!/bin/sh\nexit 101' > /usr/sbin/policy-rc.d && \
    chmod +x /usr/sbin/policy-rc.d

# Install Requirements
RUN apt-get update -qq && apt-get install -yq libc6-i386

# Add Files
ADD /run.sh /run.sh
RUN chmod 755 /*.sh
# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR
# Get AVG Installer
ADD http://download.avgfree.com/filedir/inst/avg2013flx-r3118-a6926.i386.deb /avg/

# Install AVG
RUN dpkg -i /avg/avg2013flx-r3118-a6926.i386.deb

# Start AVG Service and Update AVG Definitions
RUN /etc/init.d/avgd restart && avgupdate
RUN avgcfgctl -w UpdateVir.sched.Task.Disabled=true
RUN avgcfgctl -w Default.setup.daemonize=false

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

CMD ["/run.sh"]
