Avast Dockerfile
=============

This repository contains a **Dockerfile** of [Avast](http://www.avast.com/registration-free-antivirus.php) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/avast/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/avast/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/avast`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/avast github.com/blacktop/docker-avast`

### Usage
```bash
$ docker run -i -t blacktop/avast EICAR
```

#### Or link your own malware folder:
```
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/avast
```

#### Output:
```bash

```
### Todo
- [x] Install/Run Avast
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
