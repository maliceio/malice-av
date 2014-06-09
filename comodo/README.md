Comodo Dockerfile
=============

This repository contains a **Dockerfile** of [Comodo](https://www.comodo.com/home/internet-security/antivirus-for-linux.php) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/comodo/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/comodo/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/comodo`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/comodo github.com/blacktop/docker-comodo`
### Usage
```
$ docker run -i -t blacktop/comodo -vs /malware/EICAR
```
#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/comodo
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
- [x] Install/Run Comodo
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
