Avira Dockerfile
================

This repository contains a **Dockerfile** of [Avira](http://www.avira.com/en/index) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/avira/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies
* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation
1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/avira/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/avira`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/avira github.com/blacktop/docker-avira`

### Usage

    $ docker run -i -t blacktop/avira

#### Or link your own malware folder:

    $ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/avira

#### Output:

    Avira AntiVir Professional (ondemand scanner)
    Copyright (C) 2010 by Avira GmbH.
    All rights reserved.
    
    SAVAPI-Version: 3.1.1.8, AVE-Version: 8.3.18.22
    VDF-Version: 7.11.151.18 created 20140523
    
    AntiVir license: 2228602884
    
    Info: automatically excluding /sys/ from scan (special fs)
    Info: automatically excluding /proc/ from scan (special fs)
    Info: automatically excluding /home/quarantine/ from scan (quarantine)
    
      file: /malware/EICAR
        last modified on  date: 2014-04-15  time: 07:29:59,  size: 68 bytes
        "ALERT: Eicar-Test-Signature" ; virus ; Contains code of the Eicar-Test-Signature virus
        ALERT-URL: http://www.avira.com/en/threats?q=Eicar%2DTest%2DSignature
      no action taken
    
    ------ scan results ------
       directories: 0
     scanned files: 1
            alerts: 1
        suspicious: 0
          repaired: 0
           deleted: 0
           renamed: 0
             moved: 0
         scan time: 00:00:01
    --------------------------

### Todo
- [x] Install/Run Avira
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder

