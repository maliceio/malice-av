package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var Version string

// FileInfo json object
type FileInfo struct {
	SSDeep   string            `json:"ssdeep"`
	TRiD     []string          `json:"trid"`
	Exiftool map[string]string `json:"exiftool"`
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RunCommand runs cmd on file
func RunCommand(cmd string, path string) string {

	cmdOut, err := exec.Command(cmd, path).Output()
	assert(err)

	return string(cmdOut)
}

// ParseExiftoolOutput convert exiftool output into JSON
func ParseExiftoolOutput(exifout string) map[string]string {

	var ignoreTags = []string{
		"Directory",
		"File Name",
		"File Permissions",
		"File Modification Date/Time",
	}

	lines := strings.Split(exifout, "\n")
	datas := make(map[string]string, len(lines))

	for _, line := range lines {
		keyvalue := strings.Split(line, ":")
		if len(keyvalue) != 2 {
			continue
		}
		if !stringInSlice(strings.TrimSpace(keyvalue[0]), ignoreTags) {
			datas[strings.TrimSpace(keyvalue[0])] = strings.TrimSpace(keyvalue[1])
		}
	}

	return datas
}

// ParseSsdeepOutput convert ssdeep output into JSON
func ParseSsdeepOutput(ssdout string) string {

	// Break output into lines
	lines := strings.Split(ssdout, "\n")
	// Break second line into hash and path
	hashAndPath := strings.Split(lines[1], ",")

	return strings.TrimSpace(hashAndPath[0])
}

// ParseTRiDOutput convert trid output into JSON
func ParseTRiDOutput(tridout string) []string {

	keepLines := []string{}

	lines := strings.Split(tridout, "\n")
	lines = lines[6:]
	// fmt.Println(lines)

	for _, line := range lines {
		if len(line) != 0 {
			keepLines = append(keepLines, strings.TrimSpace(line))
		}
	}

	return keepLines
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

	path := os.Args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		assert(err)
	}

	fileInfo := FileInfo{
		SSDeep:   ParseSsdeepOutput(RunCommand("ssdeep", path)),
		TRiD:     ParseTRiDOutput(RunCommand("trid", path)),
		Exiftool: ParseExiftoolOutput(RunCommand("exiftool", path)),
	}

	fileInfoJSON, err := json.Marshal(fileInfo)
	assert(err)

	fmt.Println(string(fileInfoJSON))
}
