ClamAV Dockerfile
=============

This repository contains a **Dockerfile** of [ClamAV](http://www.clamav.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/clamav/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [debian:jessie](https://index.docker.io/_/debian/)


### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/clamav/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/clamav`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/clamav github.com/blacktop/docker-clamav`

### Usage

    docker run -i -t blacktop/clamav EICAR

#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/clamav
```
#### Output:
```bash
EICAR: 'Eicar-Test-Signature FOUND'

----------- SCAN SUMMARY -----------
Known viruses: 3324284
Engine version: 0.98.1
Scanned directories: 0
Scanned files: 1
Infected files: 1
Data scanned: 0.00 MB
Data read: 0.00 MB (ratio 0.00:1)
Time: 7.009 sec (0 m 7 s)
```
### Todo
- [x] Install/Run ClamAV
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
