malice-avira [DEAD]
===================

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/avira.svg)](https://hub.docker.com/r/malice/avira/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/avira.svg)](https://hub.docker.com/r/malice/avira/)

This repository contains a **Dockerfile** of [Avira](http://www.avira.com/en/index) for [Docker](https://www.docker.io/)'s [trusted build](https://hub.docker.com/r/malice/avira/) published to the public [DockerHub](https://hub.docker.com).

> The license expired at 2014-03-01. (Appears that FREE avira is dead) :cry:

### Dependencies

-	[ubuntu:precise (*138 MB*\)](https://hub.docker.com/_/ubuntu/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/avira/) from public [DockerHub](https://hub.docker.com): `docker pull malice/avira`

### Usage

```
docker run --rm malice/avira EICAR
```

#### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/file/:/malware:ro malice/avira

Usage: avira [OPTIONS] COMMAND [arg...]

Malice Avira AntiVirus Plugin

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

Run 'avira COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

```json
{
  "avira": {
    "infected": true,
    "result": "EICAR-Test-File (not a virus)",
    "engine": "7.90123",
    "updated": "20160227"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### Avira

| Infected | Result                        | Engine  | Updated  |
|----------|-------------------------------|---------|----------|
| true     | EICAR-Test-File (not a virus) | 7.90123 | 20160227 |

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
$ docker run --name=avira malice/avira update
```

Then to used the updated Avira container:

```bash
$ docker restart avira > /dev/null && docker exec avira scan --table EICAR
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-av/issues/new) and I'll get right on it.

### Credits

### License

MIT Copyright (c) 2016 **blacktop**

<!-- Avira Dockerfile
================

This repository contains a **Dockerfile** of [Avira](http://www.avira.com/en/index) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/avira/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies
* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation
1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/avira/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/avira`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/avira github.com/blacktop/docker-avira`

### Usage

    $ docker run -i -t blacktop/avira

#### Or link your own malware folder:

    $ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/avira

#### Output:

    Avira AntiVir Professional (ondemand scanner)
    Copyright (C) 2010 by Avira GmbH.
    All rights reserved.

    SAVAPI-Version: 3.1.1.8, AVE-Version: 8.3.18.22
    VDF-Version: 7.11.151.18 created 20140523

    AntiVir license: 2228602884

    Info: automatically excluding /sys/ from scan (special fs)
    Info: automatically excluding /proc/ from scan (special fs)
    Info: automatically excluding /home/quarantine/ from scan (quarantine)

      file: /malware/EICAR
        last modified on  date: 2014-04-15  time: 07:29:59,  size: 68 bytes
        "ALERT: Eicar-Test-Signature" ; virus ; Contains code of the Eicar-Test-Signature virus
        ALERT-URL: http://www.avira.com/en/threats?q=Eicar%2DTest%2DSignature
      no action taken

    ------ scan results ------
       directories: 0
     scanned files: 1
            alerts: 1
        suspicious: 0
          repaired: 0
           deleted: 0
           renamed: 0
             moved: 0
         scan time: 00:00:01
    --------------------------

### Todo
- [x] Install/Run Avira
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
 -->
