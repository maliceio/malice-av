malice-clamav
=============

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/clamav.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/clamav.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/clamav/latest.svg)](https://imagelayers.io/?images=malice/clamav:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/clamav/latest.svg)](https://imagelayers.io/?images=malice/clamav:latest)

This repository contains a **Dockerfile** of [ClamAV](http://www.clamav.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/clamav/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

* [gliderlabs/alpine:3.3](https://index.docker.io/_/gliderlabs/alpine/)


### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://hub.docker.com/r/malice/clamav/) from public [DockerHub](https://hub.docker.com): `docker pull malice/clamav`

### Usage

    docker run -it --rm malice/clamav EICAR

#### Or link your own malware folder:
```bash
$ docker run -it --rm -v /path/to/file/:/malware:ro malice/clamav FILE
```

#### Output JSON:
```json
{
  "infected": true,
  "result": "Eicar-Test-Signature",
  "engine": " 0.99",
  "known": " 4211363",
  "updated": "20160109"
}
```

#### Output STDOUT:
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

### To Run on OSX
 - Install [Homebrew](http://brew.sh)

```bash
$ brew install caskroom/cask/brew-cask
$ brew cask install virtualbox
$ brew install docker
$ brew install docker-machine
$ docker-machine create --driver virtualbox malice
$ eval $(docker-machine env malice)
```

### Documentation

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/clamav/
