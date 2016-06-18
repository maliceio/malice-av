#!/bin/bash

/etc/init.d/avast start > /dev/null 2>&1 && scan /malware
