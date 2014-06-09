AVG Dockerfile
=============

This repository contains a **Dockerfile** of [AVG](http://free.avg.com/us-en/homepage) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/avg/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [debian:jessie](https://index.docker.io/_/debian/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/avg/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/avg`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/avg github.com/blacktop/docker-avg`

### Usage
```bash
$ docker run -i -t blacktop/avg
```

#### Or link your own malware folder:
```
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/avg
```

#### Output:
```bash
AVG command line Anti-Virus scanner
Copyright (c) 2013 AVG Technologies CZ

Virus database version: 3722/7352
Virus database release date: Wed, 16 Apr 2014 14:30:00 +0000

EICAR  'Virus identified EICAR_Test'

Files scanned     :  1(1)
Infections found  :  1(1)
PUPs found        :  0
Files healed      :  0
Warnings reported :  0
Errors reported   :  0
```

### Todo
- [x] Install/Run AVG
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
