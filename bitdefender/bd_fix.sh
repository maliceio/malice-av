#!/bin/bash
#Written by Tal
#Fixes segfaults with bitdefender

FILE_NAME=""
FILE_CHECK=$(ls /opt/BitDefender-scanner/var/lib/scan/versions.dat.* 2>/dev/null | wc -l)

#Check if BitDefender is installed
if [[ ! -d "/opt/BitDefender-scanner" ]]
then
	echo "BitDefender not installed"
	exit 0
fi

#If symlink already exists, do nothing
if [[ -L "/opt/BitDefender-scanner/var/lib/scan/bdcore.so" ]]
then
	echo "Fix already applied"
	exit 0
fi

#Only run if user is root
uid=$(/usr/bin/id -u) && [ "$uid" = "0" ] ||
{ echo "You must be root to run $0. Try again with the command 'sudo $0'"; exit 1; }


if [[ "$FILE_CHECK" != "0" ]]
then
	#versions.dat.* file exists
	FILE_NAME=$(cat /opt/BitDefender-scanner/var/lib/scan/versions.dat.* | grep -o 'bdcore\.so\.linux[^ ]*')
else
	#versions.dat.* file doesn't exist
	if [[ "$1" != "updated" ]]
	then
		#Update antivirus database and run script again
		bdscan --update
		exec ${0} updated
	else
		#Guess the filename based on output of 'uname -m'
		FILE_NAME="bdcore.so.linux-`uname -m | sed 's/i686/x86/'`"
	fi
fi

touch /opt/BitDefender-scanner/var/lib/scan/$FILE_NAME
bdscan --update
mv -v /opt/BitDefender-scanner/var/lib/scan/bdcore.so{,.old}
ln -sv /opt/BitDefender-scanner/var/lib/scan/$FILE_NAME /opt/BitDefender-scanner/var/lib/scan/bdcore.so
chown -v bitdefender:bitdefender /opt/BitDefender-scanner/var/lib/scan/$bdcore_so
