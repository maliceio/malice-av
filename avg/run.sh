#!/bin/bash
/etc/init.d/avgd start > /dev/null 2>&1 && /usr/bin/avgscan /malware
#exec supervisord
