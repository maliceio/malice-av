malice-fsecure
===============

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/f-secure.svg)](https://hub.docker.com/r/malice/f-secure/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/f-secure.svg)](https://hub.docker.com/r/malice/f-secure/)

This repository contains a **Dockerfile** of [f-secure](https://www.f-secure.com/en/web/business_global/downloads/linux-security/latest) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/f-secure/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[debian:jessie (*125 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/f-secure/) from public [DockerHub](https://hub.docker.com): `docker pull malice/f-secure`

### Usage

```
docker run --rm malice/f-secure EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/f-secure FILE

Usage: f-secure [OPTIONS] COMMAND [arg...]

Malice F-Secure AntiVirus Plugin

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

Run 'f-secure COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "f-secure": {
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

#### F-Secure

| Infected | Result                  | Engine | Updated  |
|----------|-------------------------|--------|----------|
| true     | EICAR Test-NOT virus!!! | 2.1.2  | 20160618 |

---

### To write results to [RethinkDB](https://rethinkdb.com)

```bash
$ docker volume create --name malice
$ docker run -d -p 28015:28015 -p 8080:8080 -v malice:/data --name rethink rethinkdb
$ docker run --rm -v /path/to/malware:/malware:ro --link rethink:rethink malice/f-secure -t FILE
```

### Documentation

To update the AV run the following:

```bash
$ docker run --name=f-secure malice/f-secure update
```

Then to use the updated f-secure container:

```bash
$ docker commit f-secure malice/f-secure:updated
$ docker rm f-secure # clean up updated container
$ docker run --rm malice/f-secure:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### License

MIT Copyright (c) 2016 **blacktop**
