package main_test

import (
	"strings"
	"testing"

	"github.com/cloudflare/cfssl/log"
)

const resultString = `EVALUATION VERSION - FULLY FUNCTIONAL - FREE TO USE FOR 30 DAYS.
To purchase license, please check http://www.F-Secure.com/purchase/

F-Secure Anti-Virus CLI version 1.0  build 0060

Scan started at Mon Aug 22 02:43:50 2016
Database version: 2016-08-22_01

eicar.com.txt: Infected: EICAR_Test_File [FSE]
eicar.com.txt: Infected: EICAR-Test-File (not a virus) [Aquarius]

Scan ended at Mon Aug 22 02:43:50 2016
1 file scanned
1 file infected
`

const versionString = `EVALUATION VERSION - FULLY FUNCTIONAL - FREE TO USE FOR 30 DAYS.
To purchase license, please check http://www.F-Secure.com/purchase/

F-Secure Linux Security version 11.00 build 79

F-Secure Anti-Virus CLI Command line client version:
	F-Secure Anti-Virus CLI version 1.0  build 0060

F-Secure Anti-Virus CLI Daemon version:
	F-Secure Anti-Virus Daemon version 1.0  build 0117

Database version: 2016-09-19_01

Scanner Engine versions:
	F-Secure Corporation Hydra engine version 5.15 build 154
	F-Secure Corporation Hydra database version 2016-09-16_01

	F-Secure Corporation Aquarius engine version 1.0 build 3
	F-Secure Corporation Aquarius database version 2016-09-19_01

Portions:
Copyright (c) 1994-2010 Lua.org, PUC-Rio.
Copyright (c) Reuben Thomas 2000-2010.

For full license information on Hydra engine please see licenses-fselinux.txt in the databases folder
`

func parseFSecureVersion(versionOut string) (version string, database string) {

	lines := strings.Split(versionOut, "\n")

	for _, line := range lines {

		if strings.Contains(line, "F-Secure Linux Security version") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "F-Secure Linux Security version"))
		}

		if strings.Contains(line, "Database version:") {
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

func parseFSecureOutput(fsecureout string) (map[string]string, error) {

	// root@70bc84b1553c:/malware# fsav --virus-action1=none eicar.com.txt
	// EVALUATION VERSION - FULLY FUNCTIONAL - FREE TO USE FOR 30 DAYS.
	// To purchase license, please check http://www.F-Secure.com/purchase/
	//
	// F-Secure Anti-Virus CLI version 1.0  build 0060
	//
	// Scan started at Mon Aug 22 02:43:50 2016
	// Database version: 2016-08-22_01
	//
	// eicar.com.txt: Infected: EICAR_Test_File [FSE]
	// eicar.com.txt: Infected: EICAR-Test-File (not a virus) [Aquarius]
	//
	// Scan ended at Mon Aug 22 02:43:50 2016
	// 1 file scanned
	// 1 file infected

	// log.Debugln(fsecureout)

	fsecure := make(map[string]string)

	lines := strings.Split(fsecureout, "\n")

	for _, line := range lines {
		if strings.Contains(line, "Infected:") && strings.Contains(line, "[FSE]") {
			parts := strings.Split(line, "Infected:")
			fsecure["fse"] = strings.TrimSuffix(parts[1], "[FSE]")
			continue
		}
		if strings.Contains(line, "Infected:") && strings.Contains(line, "[Aquarius]") {

			parts := strings.Split(line, "Infected:")
			fsecure["aquarius"] = strings.TrimSuffix(parts[1], "[Aquarius]")
		}
	}

	return fsecure, nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseResult(t *testing.T) {

	results, err := parseFSecureOutput(resultString)

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseVersion(t *testing.T) {

	version, database := parseFSecureVersion(versionString)

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
