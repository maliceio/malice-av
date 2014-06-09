Bitdefender Dockerfile
================

This repository contains a **Dockerfile** of [Bitdefender](http://www.bitdefender.com/business/antivirus-for-unices.html) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/bitdefender/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies
* [ubuntu:latest](https://index.docker.io/_/ubuntu/)

### Installation
1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/bitdefender/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/bitdefender`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/bitdefender github.com/blacktop/docker-bitdefender`

### Usage
```bash
$ docker run -i -t blacktop/bitdefender
```
#### Or link your own malware folder:
```bash
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/bitdefender
```
#### Output:
```bash
BitDefender Antivirus Scanner for Unices v7.90123 Linux-amd64
Copyright (C) 1996-2009 BitDefender. All rights reserved.
Trial key found. 30 days remaining.

Infected file action: ignore
Suspected file action: ignore
Loading plugins, please wait
Plugins loaded.

/malware/EICAR  infected: 'EICAR-Test-File (not a virus)'


Results:
Folders: 0
Files: 1
Packed: 0
Archives: 0
Infected files: 1
Suspect files: 0
Warnings: 0
Identified viruses: 1
I/O errors: 0
```

### Todo
- [x] Install/Run Bitdefender
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
