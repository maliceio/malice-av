# malice-fprot

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/fprot.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/fprot.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/fprot/latest.svg)](https://imagelayers.io/?images=malice/fprot:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/fprot/latest.svg)](https://imagelayers.io/?images=malice/fprot:latest)

This repository contains a **Dockerfile** of [fprot](http://www.fprot.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/fprot/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

* [debian:jessie (*125 MB*)](https://index.docker.io/_/debian/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://hub.docker.com/r/malice/fprot/) from public [DockerHub](https://hub.docker.com): `docker pull malice/fprot`

### Usage

    docker run -it --rm malice/fprot EICAR

#### Or link your own malware folder:
```bash
$ docker run -it --rm -v /path/to/file/:/malware:ro malice/fprot FILE
```
#### Output JSON:
```json
{
  "f-prot": {
    "infected": true,
    "result": "EICAR_Test_File (exact)",
    "engine": "4.6.5.141",
    "updated": "201601110435"
  }
}
```
#### Output STDOUT:
```bash
F-PROT Antivirus CLS version 6.7.10.6267, 64bit (built: 2012-03-27T11-39-07)


FRISK Software International (C) Copyright 1989-2011
Engine version:   4.6.5.141
Arguments:        -r EICAR
Virus signatures: 201601070641
                  (/opt/f-prot/antivir.def)

[Found virus] <EICAR_Test_File (exact)> 	EICAR
Scanning: /

Results:

Files: 1
Skipped files: 0
MBR/boot sectors checked: 0
Objects scanned: 1
Infected objects: 1
Infected files: 1
Files with errors: 0
Disinfected: 0

Running time: 00:00
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

[hub]: https://hub.docker.com/r/malice/fprot/
