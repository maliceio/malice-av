malice-avg
==========

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/avg.svg)](https://hub.docker.com/r/malice/avg/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/avg.svg)](https://hub.docker.com/r/malice/avg/)

This repository contains a **Dockerfile** of [avg](http://www.avg.net/lang/en/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/avg/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[debian:jessie (*125 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/avg/) from public [DockerHub](https://hub.docker.com): `docker pull malice/avg`

### Usage

```
docker run --rm malice/avg EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/avg FILE

Usage: avg [OPTIONS] COMMAND [arg...]

Malice AVG AntiVirus Plugin

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
    "updated": "20160213"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### AVG

| Infected | Result                      | Engine    | Updated  |
|----------|-----------------------------|-----------|----------|
| true     | Virus identified EICAR_Test | 13.0.3114 | 20160213 |

---

### To write results to [ElasticSearch](https://www.elastic.co/products/elasticsearch)

```bash
$ docker volume create --name malice
$ docker run -d -p 9200:9200 -v malice:/data --name elastic elasticsearch
$ docker run --rm -v /path/to/malware:/malware:ro --link elastic malice/avg -t FILE
```

### Documentation

To update the AV run the following:

```bash
$ docker run --name=avg malice/avg update
```

Then to use the updated AVG container:

```bash
$ docker commit avg malice/avg:updated
$ docker rm avg # clean up updated container
$ docker run --rm malice/avg:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### License

MIT Copyright (c) 2016 **blacktop**
