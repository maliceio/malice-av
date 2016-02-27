# malice-bitdefender

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/bitdefender.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/bitdefender.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/bitdefender/latest.svg)](https://imagelayers.io/?images=malice/bitdefender:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/bitdefender/latest.svg)](https://imagelayers.io/?images=malice/bitdefender:latest)

This repository contains a **Dockerfile** of [Bitdefender](http://www.bitdefender.com/business/antivirus-for-unices.html) for [Docker](https://www.docker.io/)'s [trusted build][hub] published to the public [DockerHub](https://hub.docker.com).

### Dependencies

* [ubuntu:precise (*138 MB*)](https://hub.docker.com/_/ubuntu/)

### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build][hub] from public [DockerHub](https://hub.docker.com): `docker pull malice/bitdefender`

### Usage

    docker run -it --rm malice/bitdefender EICAR

#### Or link your own malware folder:
```bash
$ docker run -it --rm -v /path/to/file/:/malware:ro malice/bitdefender

Usage: bitdefender [OPTIONS] COMMAND [arg...]

Malice Bitdefender AntiVirus Plugin

Version: v0.1.0, BuildTime: 20160227

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

Run 'bitdefender COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:
```json
{
  "bitdefender": {
    "infected": true,
    "result": "Malware",
    "engine": "1.1",
    "updated": "20160227"
  }
}
```
### Sample Output STDOUT (Markdown Table):
---
#### Bitdefender
| Infected | Result  | Engine | Updated  |
| -------- | ------- | ------ | -------- |
| true     | Malware | 1.1    | 20160227 |
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
$ docker run --name=bitdefender malice/bitdefender update
```
Then to used the updated Bitdefender container:
```bash
$ docker restart bitdefender > /dev/null && docker exec bitdefender scan --table EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/bitdefender/
