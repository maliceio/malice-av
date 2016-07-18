FROM debian:jessie
MAINTAINER blacktop, https://github.com/blacktop

# Install Requirements
RUN apt-get update && apt-get install -y wget \
  && wget -qO - http://repo.drweb.com/drweb/drweb.key | apt-key add - \
  && echo "deb http://repo.drweb.com/drweb/debian 10.0.0 non-free" >> /etc/apt/sources.list \
  && apt-get update && apt-get install -y drweb-workstations
