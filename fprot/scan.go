package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// FPROT json object
type FPROT struct {
	Infected bool   `json:"infected"`
	Result   string `json:"result"`
	Engine   string `json:"engine"`
	Updated  string `json:"updated"`
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
func RunCommand(cmd string, args ...string) string {

	cmdOut, err := exec.Command(cmd, args...).Output()
	if len(cmdOut) == 0 {
		assert(err)
	}

	return string(cmdOut)
}

// ParseFprotOutput convert fprot output into FPROT struct
func ParseFprotOutput(fprotout string) FPROT {

	fprot := FPROT{Infected: false}
	colonSeparated := []string{}

	lines := strings.Split(fprotout, "\n")
	// Extract Virus string and extract colon separated lines into an slice
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
			if strings.Contains(line, "[Found virus]") {
				result := extractVirusName(line)
				if len(result) != 0 {
					fprot.Result = result
					fprot.Infected = true
				} else {
					fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
					os.Exit(2)
				}
			}
		}
	}
	// fmt.Println(lines)

	// Extract FPROT Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := strings.Split(line, ":")
				if len(keyvalue) != 0 {
					switch {
					case strings.Contains(keyvalue[0], "Virus signatures"):
						fprot.Updated = strings.TrimSpace(keyvalue[1])
					case strings.Contains(line, "Engine version"):
						fprot.Engine = strings.TrimSpace(keyvalue[1])
					}
				}
			}
		}
	} else {
		fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
		os.Exit(2)
	}

	return fprot
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	r := regexp.MustCompile(`<(.+)>`)
	res := r.FindStringSubmatch(line)
	if len(res) != 2 {
		return ""
	}
	return res[1]
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("[ERROR] Missing input file.")
		os.Exit(2)
	}
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("Version: ", Version)
		fmt.Println("BuildTime: ", BuildTime)
		os.Exit(0)
	}

	path := os.Args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		assert(err)
	}

	fprotOutput := RunCommand("/usr/local/bin/fpscan", "-r", path)
	// fmt.Println(ParseFprotOutput(clamOutput))

	fprotJSON, err := json.Marshal(ParseFprotOutput(fprotOutput))
	assert(err)

	fmt.Println(string(fprotJSON))
}
