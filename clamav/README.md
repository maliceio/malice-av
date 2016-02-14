# malice-clamav

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
$ docker run -it --rm -v /path/to/file/:/malware:ro malice/clamav

Usage: clamav [OPTIONS] COMMAND [arg...]

Malice ClamAV Plugin

Version: v0.1.0, BuildTime: 20160214

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --table, -t	output as Markdown table
  --post, -p	POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x	proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --help, -h	show help
  --version, -v	print the version

Commands:
  update	Update virus definitions
  help		Shows a list of commands or help for one command

Run 'clamav COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:
```json
{
  "clamav": {
    "infected": true,
    "result": "Eicar-Test-Signature",
    "engine": "0.99",
    "known": "4213581",
    "updated": "20160213"
  }
}
```
### Sample Output STDOUT (Markdown Table):
---
#### ClamAV
| Infected | Result               | Engine | Updated  |
| -------- | -------------------- | ------ | -------- |
| true     | Eicar-Test-Signature | 0.99   | 20160213 |
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
To update the AV run the following:
```bash
$ docker run --name=clamav malice/clamav update
```
Then to used the updated ClamAV container:
```bash
$ docker restart clamav > /dev/null && docker exec clamav scan --table EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/clamav/
