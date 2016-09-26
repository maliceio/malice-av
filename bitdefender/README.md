malice-bitdefender
==================

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/bitdefender.svg)](https://hub.docker.com/r/malice/bitdefender/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/bitdefender.svg)](https://hub.docker.com/r/malice/bitdefender/)

This repository contains a **Dockerfile** of [Bitdefender](http://www.bitdefender.com/business/antivirus-for-unices.html) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/bitdefender/) published to the public [DockerHub](https://hub.docker.com).

### Dependencies

-	[ubuntu:precise (*138 MB*\)](https://hub.docker.com/_/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/bitdefender/) from public [DockerHub](https://hub.docker.com): `docker pull malice/bitdefender`

### Usage

```
docker run --rm malice/bitdefender EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/bitdefender FILE

Usage: bitdefender [OPTIONS] COMMAND [arg...]

Malice Bitdefender AntiVirus Plugin

Version: v0.1.0, BuildTime: 20160227

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V         verbose output
  --table, -t           output as Markdown table
  --post, -p            POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x           proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --rethinkdb value     rethinkdb address for Malice to store results [$MALICE_RETHINKDB]
  --help, -h            show help
  --version, -v         print the version

Commands:
  update        Update virus definitions
  help          Shows a list of commands or help for one command

Run 'bitdefender COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "bitdefender": {
    "infected": true,
    "result": "EICAR-Test-File (not a virus)",
    "engine": "7.90123",
    "updated": "20160227"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### Bitdefender

| Infected | Result                        | Engine  | Updated  |
|----------|-------------------------------|---------|----------|
| true     | EICAR-Test-File (not a virus) | 7.90123 | 20160227 |

---

### To write results to [ElasticSearch](https://www.elastic.co/products/elasticsearch)

```bash
$ docker volume create --name malice
$ docker run -d -p 9200:9200 -v malice:/data --name elastic elasticsearch
$ docker run --rm -v /path/to/malware:/malware:ro --link elastic malice/bitdefender -t FILE
```

### Documentation

To update the AV run the following:

```bash
$ docker run --name=bitdefender malice/bitdefender update
```

Then to use the updated Bitdefender container:

```bash
$ docker commit bitdefender malice/bitdefender:updated
$ docker rm bitdefender # clean up updated container
$ docker run --rm malice/bitdefender:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### License

MIT Copyright (c) 2016 **blacktop**
