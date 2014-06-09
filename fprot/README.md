F-PROT Dockerfile
=============

This repository contains a **Dockerfile** of [F-PROT](http://www.f-prot.com/products/home_use/linux/) for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/blacktop/fprot/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [debian:jessie](https://index.docker.io/_/debian/)


### Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://index.docker.io/u/blacktop/fprot/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/fprot`

#### Alternatively, build an image from Dockerfile
`docker build -t blacktop/fprot github.com/blacktop/docker-fprot`

### Usage
```bash
$ docker run -i -t blacktop/fprot -r EICAR
```
#### Or link your own malware folder:
```
$ docker run -i -t -v /path/to/malware/:/malware:ro blacktop/fprot -r /malware
```
#### Output:
```bash
F-PROT Antivirus CLS version 6.7.10.6267, 32bit (built: 2012-03-27T12-34-14)


FRISK Software International (C) Copyright 1989-2011
Engine version:   4.6.5.141
Arguments:        -r EICAR
Virus signatures: 201404190947
                  (/opt/f-prot/antivir.def)

[Found virus] '<EICAR_Test_File (exact)> 	EICAR'
Scanning: /

Results:

Files: 1
Skipped files: 0
MBR/boot sectors checked: 0
Objects scanned: 1
Infected objects: 1
Infected files: 1
Files with errors: 0
Disinfected: 0

Running time: 00:01
```
### Todo
- [x] Install/Run F-PROT
- [ ] Start Daemon and watch folder with supervisord
- [ ] Have container take a URL as input and download/scan file
- [ ] Output Scan Results as formated JSON
- [ ] Attach a Volume that will hold malware for a host's tmp folder
