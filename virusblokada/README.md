
version: /opt/vba/vbacl --version

    - libc6-i386
    - libstdc++6

    dpkg --add-architecture i386

http://anti-virus.by/pub/vbacl-linux-3.12.26.4.tar.gz

update: vbacl --update

install: setup.sh install
