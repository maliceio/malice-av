package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// AVG json object
type AVG struct {
	Results ResultsData `json:"avg"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected"`
	Result   string `json:"result"`
	Engine   string `json:"engine"`
	Database string `json:"database"`
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

// ParseAVGOutput convert avg output into ResultsData struct
func ParseAVGOutput(avgout string, path string) ResultsData {

	avg := ResultsData{
		Infected: false,
		Engine:   getAvgVersion(),
	}
	colonSeparated := []string{}

	lines := strings.Split(avgout, "\n")
	// Extract Virus string and extract colon separated lines into an slice
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, ":") {
				colonSeparated = append(colonSeparated, line)
			}
			if strings.Contains(line, path) {
				pathVirusString := strings.Split(line, "  ")
				avg.Result = strings.TrimSpace(pathVirusString[1])
			}
		}
	}
	// fmt.Println(lines)

	// Extract AVG Details from scan output
	if len(colonSeparated) != 0 {
		for _, line := range colonSeparated {
			if len(line) != 0 {
				keyvalue := strings.Split(line, ":")
				if len(keyvalue) != 0 {
					switch {
					case strings.Contains(line, "Virus database version"):
						avg.Database = strings.TrimSpace(keyvalue[1])
					case strings.Contains(line, "Virus database release date"):
						date := strings.TrimSpace(strings.TrimPrefix(line, "Virus database release date:"))
						avg.Updated = parseUpdatedDate(date)
					case strings.Contains(line, "Infections found"):
						if strings.Contains(keyvalue[1], "1") {
							avg.Infected = true
						}
					}
				}
			}
		}
	} else {
		fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
		fmt.Printf("[ERROR] AVG output was: \n%s", avgout)
		os.Exit(2)
	}

	return avg
}

// Get Anti-Virus scanner version
func getAvgVersion() string {
	versionOut := RunCommand("/usr/bin/avgscan", "-v")
	lines := strings.Split(versionOut, "\n")
	for _, line := range lines {
		if len(line) != 0 {
			keyvalue := strings.Split(line, ":")
			if len(keyvalue) != 0 {
				if strings.Contains(keyvalue[0], "Anti-Virus scanner version") {
					return strings.TrimSpace(keyvalue[1])
				}
			}
		}
	}
	return ""
}

func parseUpdatedDate(date string) string {
	layout := "Mon, 02 Jan 2006 15:04:05 +0000"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func updateAV() {
	fmt.Println("Updating AVG...")
	// AVG needs to have the daemon started first
	exec.Command("/etc/init.d/avgd", "start").Output()

	fmt.Println(RunCommand("avgupdate"))
}

func printMarkDownTable(avg AVG) {

	fmt.Println("#### AVG")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": avg.Results.Infected,
		"Result":   avg.Results.Result,
		"Engine":   avg.Results.Engine,
		"Updated":  avg.Results.Updated,
	})
	table.Markdown = true
	table.Print()
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
	app.Name = "avg"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice AVG AntiVirus Plugin"
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
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}

		// AVG needs to have the daemon started first
		exec.Command("/etc/init.d/avgd", "start").Output()

		avg := AVG{
			Results: ParseAVGOutput(RunCommand("/usr/bin/avgscan", path), path),
		}

		if c.Bool("table") {
			printMarkDownTable(avg)
		} else {
			avgJSON, err := json.Marshal(avg)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(avgJSON).
					End(printStatus)
			}
			fmt.Println(string(avgJSON))
		}
	}

	err := app.Run(os.Args)
	assert(err)
}
