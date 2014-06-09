Dr.Web Dockerfile
=============

This repository contains a **Dockerfile** of [Dr.Web](https://download.drweb.com/linux/?lng=en) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/drweb/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/drweb/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/drweb`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/drweb github.com/blacktop/docker-drweb`
### Usage
```
$ docker run -i -t blacktop/drweb EICAR
```
#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/drweb
```
#### Output:
```bash

```
### Todo
- [x] Install/Run Dr.Web
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder