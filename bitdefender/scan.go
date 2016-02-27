package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/parnurzeal/gorequest"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// Bitdefender json object
type Bitdefender struct {
	Results ResultsData `json:"bitdefender"`
}

// ResultsData json object
type ResultsData struct {
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

// ParseBitdefenderOutput convert bitdefender output into ResultsData struct
func ParseBitdefenderOutput(bitdefenderout string) ResultsData {

	bitdefender := ResultsData{Infected: false, Updated: BuildTime}
	// EXAMPLE OUTPUT:
	// BitDefender Antivirus Scanner for Unices v7.90123 Linux-amd64
	// Copyright (C) 1996-2009 BitDefender. All rights reserved.
	// Trial key found. 30 days remaining.
	//
	// Infected file action: ignore
	// Suspected file action: ignore
	// Loading plugins, please wait
	// Plugins loaded.
	//
	// /malware/EICAR  infected: EICAR-Test-File (not a virus)
	//
	//
	// Results:
	// Folders: 0
	// Files: 1
	// Packed: 0
	// Archives: 0
	// Infected files: 1
	// Suspect files: 0
	// Warnings: 0
	// Identified viruses: 1
	// I/O errors: 0
	lines := strings.Split(bitdefenderout, "\n")

	// Extract Virus string
	for _, line := range lines {
		if len(line) != 0 {
			switch {
			case strings.Contains(line, "infected:"):
				result := extractVirusName(line)
				if len(result) != 0 {
					bitdefender.Result = result
					bitdefender.Infected = true
				} else {
					fmt.Println("[ERROR] Virus name extracted was empty: ", result)
					os.Exit(2)
				}
			case strings.Contains(line, "Unices v"):
				words := strings.Fields(line)
				for _, word := range words {
					if strings.HasPrefix(word, "v") {
						bitdefender.Engine = strings.TrimPrefix(word, "v")
					}
				}
			}
		}
	}
	if found, updated := checkUpdatedDate(); found {
		bitdefender.Updated = updated
	}
	return bitdefender
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	keyvalue := strings.Split(line, "infected:")
	return strings.TrimSpace(keyvalue[1])
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func checkUpdatedDate() (bool, string) {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return false, ""
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return true, string(updated)

}

func printMarkDownTable(bitdefender Bitdefender) {

	fmt.Println("#### Bitdefender")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": bitdefender.Results.Infected,
		"Result":   bitdefender.Results.Result,
		"Engine":   bitdefender.Results.Engine,
		"Updated":  bitdefender.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

func updateAV() {
	fmt.Println("Updating Bitdefender...")
	fmt.Println(RunCommand("bdscan", "--update"))
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	assert(err)
}

var appHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

{{.Usage}}

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

func main() {
	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()
	app.Name = "bitdefender"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Bitdefender AntiVirus Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "post, p",
			Usage:  "POST results to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update virus definitions",
			Action: func(c *cli.Context) {
				updateAV()
			},
		},
	}
	app.Action = func(c *cli.Context) {
		path, err := filepath.Abs(c.Args().First())
		assert(err)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}

		bitdefender := Bitdefender{
			Results: ParseBitdefenderOutput(RunCommand("bdscan", path)),
		}

		if c.Bool("table") {
			printMarkDownTable(bitdefender)
		} else {
			bitdefenderJSON, err := json.Marshal(bitdefender)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(bitdefenderJSON).
					End(printStatus)
			}
			fmt.Println(string(bitdefenderJSON))
		}
	}

	err := app.Run(os.Args)
	assert(err)
}
