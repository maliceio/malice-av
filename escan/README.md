eScan Dockerfile
================

This repository contains a **Dockerfile** of [eScan](http://www.escanav.com/english/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/escan/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [ubuntu:latest](https://index.docker.io/_/ubuntu/)


### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/escan/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/escan`

#### Alternatively, build an image from Dockerfile
```
$ docker build -t blacktop/escan github.com/blacktop/docker-av/tree/escan
```
### Usage
```
$ docker run -i -t blacktop/escan
```
#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/escan
```
#### Output:
```

```
### Todo
- [x] Install/Run eScan
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
