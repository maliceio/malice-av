F-Secure Dockerfile
=============

This repository contains a **Dockerfile** of [F-Secure](http://www.f-secure.com/en/web/business_global/support/downloads/-/carousel/view/83) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/fsecure/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/fsecure/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/fsecure`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/fsecure github.com/blacktop/docker-fsecure`
### Usage
```
$ docker run -i -t blacktop/fsecure -vs /malware/EICAR
```
#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/fsecure
```
#### Output:
```
-----== Scan Start ==-----
/malware/EICAR ---> Found Virus, Malware Name is ApplicUnwnt
-----== Scan End ==-----
Number of Scanned Files: 1
Number of Found Viruses: 1
```
### Todo
- [x] Install/Run F-Secure
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder