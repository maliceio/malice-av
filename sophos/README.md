malice-sophos
=============

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/sophos.svg)](https://hub.docker.com/r/malice/sophos/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/sophos.svg)](https://hub.docker.com/r/malice/sophos/)

This repository contains a **Dockerfile** of [Sophos](https://www.sophos.com/en-us/products/free-tools/sophos-antivirus-for-linux.aspx) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/sophos/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[debian:jessie (*125 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/sophos/) from public [DockerHub](https://hub.docker.com): `docker pull malice/sophos`

### Usage

```
docker run --rm malice/sophos EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/sophos FILE

Usage: sophos [OPTIONS] COMMAND [arg...]

Malice Sophos AntiVirus Plugin

Version: v0.1.0, BuildTime: 20160618

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

Run 'sophos COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "sophos": {
    "infected": true,
    "result": "EICAR Test-NOT virus!!!",
    "engine": "2.1.2",
    "database": "16061703",
    "updated": "20160618"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### Sophos

| Infected | Result                  | Engine | Updated  |
|----------|-------------------------|--------|----------|
| true     | EICAR Test-NOT virus!!! | 2.1.2  | 20160618 |

---

### To write results to [RethinkDB](https://rethinkdb.com)

```bash
$ docker volume create --name malice
$ docker run -d -p 28015:28015 -p 8080:8080 -v malice:/data --name rethink rethinkdb
$ docker run --rm -v /path/to/malware:/malware:ro --link rethink:rethink malice/sophos -t FILE
```

### Documentation

To update the AV run the following:

```bash
$ docker run --name=sophos malice/sophos update
```

Then to use the updated sophos container:

```bash
$ docker commit sophos malice/sophos:updated
$ docker rm sophos # clean up updated container
$ docker run --rm malice/sophos:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### License

MIT Copyright (c) 2016 **blacktop**
