malice-comodo
=============

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)[![Docker Stars](https://img.shields.io/docker/stars/malice/comodo.svg)](https://hub.docker.com/r/malice/comodo/)[![Docker Pulls](https://img.shields.io/docker/pulls/malice/comodo.svg)](https://hub.docker.com/r/malice/comodo/)

This repository contains a **Dockerfile** of [Comodo](https://www.comodo.com/home/internet-security/antivirus-for-linux.php) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/comodo/) published to the public [DockerHub](https://hub.docker.com).

### Dependencies

-	[ubuntu:precise (*138 MB*\)](https://hub.docker.com/_/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/comodo/) from public [DockerHub](https://hub.docker.com): `docker pull malice/comodo`

### Usage

```
docker run --rm malice/comodo EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/comodo FILE

Usage: comodo [OPTIONS] COMMAND [arg...]

Malice Comodo AntiVirus Plugin

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

Run 'comodo COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "comodo": {
    "infected": true,
    "result": "Malware",
    "engine": "1.1",
    "updated": "20160227"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### Comodo

| Infected | Result  | Engine | Updated  |
|----------|---------|--------|----------|
| true     | Malware | 1.1    | 20160227 |

---

### To Run on OSX

-	Install [Homebrew](http://brew.sh)

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
$ docker run --name=comodo malice/comodo update
```

Then to use the updated Comodo container:

```bash
$ docker commit comodo malice/comodo:updated
$ docker rm comodo # clean up updated container
$ docker run --rm malice/comodo:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License

MIT Copyright (c) 2016 **blacktop**
