package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// ClamAV json object
type ClamAV struct {
	Infected bool   `json:"infected"`
	Result   string `json:"result"`
	Engine   string `json:"engine"`
	Known    string `json:"known"`
	Update   string `json:"update"`
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// RunCommand runs cmd on file
func RunCommand(cmd string, path string) string {

	cmdOut, err := exec.Command(cmd, path).Output()
	assert(err)

	return string(cmdOut)
}

// ParseClamAvOutput convert clamav output into ClamAV struct
func ParseClamAvOutput(clamout string) ClamAV {

	clamAV := ClamAV{}

	lines := strings.Split(clamout, "\n")
	// fmt.Println(lines)
	// Extract AV Scan Result
	result := lines[0]
	if len(result) != 0 {
		pathAndResult := strings.Split(result, ":")
		if strings.Contains(pathAndResult[1], "OK") {
			clamAV.Infected = false
		} else {
			clamAV.Infected = true
			clamAV.Result = pathAndResult[1]
		}
	} else {
		fmt.Println("[ERROR] empty scan result: ", result)
		os.Exit(2)
	}
	// Extract Clam Details from SCAN SUMMARY
	for _, line := range lines[1:] {
		if len(line) != 0 {
			keyvalue := strings.Split(line, ":")
			if len(keyvalue) != 0 {
				switch {
				case strings.Contains(keyvalue[0], "Known viruses"):
					clamAV.Known = keyvalue[1]
				case strings.Contains(line, "Engine version"):
					clamAV.Engine = keyvalue[1]
				}
			}
		}
	}

	clamAV.Update = BuildTime

	return clamAV
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("[ERROR] Missing input file.")
		os.Exit(2)
	}
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(Version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--build" {
		fmt.Println(BuildTime)
		os.Exit(0)
	}

	path := os.Args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		assert(err)
	}

	clamOutput := RunCommand("/usr/bin/clamscan", path)
	// fmt.Println(ParseClamAvOutput(clamOutput))

	clamavJSON, err := json.Marshal(ParseClamAvOutput(clamOutput))
	assert(err)

	fmt.Println(string(clamavJSON))
}
