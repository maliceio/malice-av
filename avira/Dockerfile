FROM ubuntu:latest
MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get -q update && apt-get install -yq libc6-i386

# Add Files
ADD /run.sh /run.sh
RUN chmod 755 /run.sh
# Add EICAR Test Virus File to malware folder
RUN mkdir malware && echo "X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*" > /malware/EICAR
ADD /unattended.inf /unattended.inf
#ADD /supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir /home/quarantine/ #&& mkdir -p /var/log/supervisor

# Download Avira
ADD http://premium.avira-update.com/package/wks_avira/unix/en/pers/antivir_workstation-pers.tar.gz /
RUN tar -zxvf /antivir_workstation-pers.tar.gz

# Install Avira
RUN /antivir-workstation-pers-3.1.3.5-0/install --inf=/unattended.inf
ADD http://personal.avira-update.com/package/peclkey/win32/int/hbedv.key /usr/lib/AntiVir/guard/avira.key
# Update Avira
RUN /usr/lib/AntiVir/guard/avupdate-guard --product=Guard

# Try to reduce size of container.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /malware

CMD ["/run.sh"]
