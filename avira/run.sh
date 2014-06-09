#!/bin/bash
echo "Scanning Malware..."
/usr/lib/AntiVir/guard/avguard restart > /dev/null 2>&1
/usr/lib/AntiVir/guard/avscan -s --scan-in-archive=yes \
                                 --scan-mode=all \
                                 --heur-level=3 \
                                 --alert-action=none \
                                 --heur-macro=yes \
                                 --batch /malware
# exec supervisord
