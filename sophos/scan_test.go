package main_test

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

const resultString = `SAVScan virus detection utility
Version 5.21.0 [Linux/AMD64]
Virus data version 5.27, April 2016
Includes detection for 11283995 viruses, Trojans and worms
Copyright (c) 1989-2016 Sophos Limited. All rights reserved.

System time 03:48:15, System date 22 August 2016
Command line qualifiers are: -f

Full Scanning

>>> Virus 'EICAR-AV-Test' found in file EICAR

1 file scanned in 4 seconds.
1 virus was discovered.
1 file out of 1 was infected.
If you need further advice regarding any detections please visit our
Threat Center at: http://www.sophos.com/en-us/threat-center.aspx
End of Scan.
`

const versionString = `SAVScan virus detection utility
Copyright (c) 1989-2016 Sophos Limited. All rights reserved.

System time 03:41:05, System date 22 August 2016

Product version           : 5.21.0
Engine version            : 3.64.0
Virus data version        : 5.27
User interface version    : 2.03.064
Platform                  : Linux/AMD64
Released                  : 26 April 2016
Total viruses (with IDEs) : 11283995
`

func parseSophosVersion(versionOut string) (version string, database string) {

	lines := strings.Split(versionOut, "\n")

	for _, line := range lines {

		if strings.Contains(line, "Product version") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				version = strings.TrimSpace(parts[1])
			} else {
				log.Error("Umm... ", parts)
			}
		}

		if strings.Contains(line, "Virus data version") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				database = strings.TrimSpace(parts[1])
				break
			} else {
				log.Error("Umm... ", parts)
			}
		}

	}

	return
}

func parseSophosOutput(sophosout string) (string, error) {

	lines := strings.Split(sophosout, "\n")

	for _, line := range lines {
		if strings.Contains(line, ">>> Virus") && strings.Contains(line, "found in file") {
			parts := strings.Split(line, "'")
			fmt.Println(parts)
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseResult(t *testing.T) {

	results, err := parseSophosOutput(resultString)

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseVersion(t *testing.T) {

	version, database := parseSophosVersion(versionString)

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
