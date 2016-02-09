# malice-avg

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/avg.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/avg.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/avg/latest.svg)](https://imagelayers.io/?images=malice/avg:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/avg/latest.svg)](https://imagelayers.io/?images=malice/avg:latest)

This repository contains a **Dockerfile** of [avg](http://www.avg.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/avg/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

* [debian:jessie (*125 MB*)](https://index.docker.io/_/debian/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://hub.docker.com/r/malice/avg/) from public [DockerHub](https://hub.docker.com): `docker pull malice/avg`

### Usage

    docker run -it --rm malice/avg EICAR

#### Or link your own malware folder:
```bash
$ docker run -it --rm -v /path/to/file/:/malware:ro malice/avg

Usage: avg [OPTIONS] COMMAND [arg...]

Malice AVG AntiVirus Plugin

Version: v0.1.0, BuildTime: 20160209

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --table, -t	output as Markdown table
  --post, -p	POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x	proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --help, -h	show help
  --version, -v	print the version

Commands:
  help	Shows a list of commands or help for one command

Run 'avg COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:
```json
{
  "avg": {
    "infected": true,
    "result": "Virus identified EICAR_Test",
    "engine": "13.0.3114",
    "database": "4477/11588",
    "updated": "Tue, 09 Feb 2016 00:27:00 +0000"
  }
}
```
### Sample Output STDOUT (Markdown Table):
---
#### AVG
| Infected | Result                      | Version   | Updated                         |
| -------- | --------------------------- | --------- | ------------------------------- |
| true     | Virus identified EICAR_Test | 13.0.3114 | Tue, 09 Feb 2016 00:27:00 +0000 |
---
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

[hub]: https://hub.docker.com/r/malice/avg/
