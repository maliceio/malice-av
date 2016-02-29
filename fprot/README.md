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

    docker run --rm malice/fprot EICAR

#### Or link your own malware folder:
```bash
$ docker run --rm -v /path/to/file/:/malware:ro malice/fprot

Usage: fprot [OPTIONS] COMMAND [arg...]

Malice F-PROT AntiVirus Plugin

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

Run 'fprot COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:
```json
{
  "f-prot": {
    "infected": true,
    "result": "EICAR_Test_File (exact)",
    "engine": "4.6.5.141",
    "updated": "20160213"
  }
}
```
### Sample Output STDOUT (Markdown Table):
---
#### F-PROT
| Infected | Result                  | Engine    | Updated    |
| -------- | ----------------------- | --------- | ---------- |
| true     | EICAR_Test_File (exact) | 4.6.5.141 | 20160213   |
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
$ docker run --name=fprot malice/fprot update
```
Then to use the updated F-PROT container:
```bash
$ docker commit fprot malice/fprot:updated
$ docker rm fprot # clean up updated container
$ docker run --rm malice/fprot:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/fprot/
