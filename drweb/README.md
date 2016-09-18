malice-drweb
=============

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/drweb.svg)](https://hub.docker.com/r/malice/drweb/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/drweb.svg)](https://hub.docker.com/r/malice/drweb/)

This repository contains a **Dockerfile** of [Dr.Web](https://www.drweb.com) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/drweb/) published to the public [DockerHub](https://index.docker.io/).

### Dependencies

-	[debian:jessie (*125 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/drweb/) from public [DockerHub](https://hub.docker.com): `docker pull malice/drweb`

### Usage

```
docker run --rm malice/drweb EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/drweb FILE

Usage: drweb [OPTIONS] COMMAND [arg...]

Malice Dr.Web Plugin

Version: v0.1.0, BuildTime: 20160214

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

Run 'drweb COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "drweb": {
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

#### Dr.Web

| Infected | Result               | Engine | Updated  |
|----------|----------------------|--------|----------|
| true     | Eicar-Test-Signature | 0.99   | 20160213 |

---

### To write results to [RethinkDB](https://rethinkdb.com)

```bash
$ docker volume create --name malice
$ docker run -d -p 28015:28015 -p 8080:8080 -v malice:/data --name rethink rethinkdb
$ docker run --rm -v /path/to/malware:/malware:ro --link rethink:rethink malice/drweb -t FILE
```

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
$ docker run --name=drweb malice/drweb update
```

Then to use the updated Dr.Web container:

```bash
$ docker commit drweb malice/drweb:updated
$ docker rm drweb # clean up updated container
$ docker run --rm malice/drweb:updated EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License

MIT Copyright (c) 2016 **blacktop**
