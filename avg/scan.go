package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/fatih/structs"
	"github.com/maliceio/go-plugin-utils/database/elasticsearch"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "avg"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"avast" gorethink:"avg"`
}

// AVG json object
type AVG struct {
	Results ResultsData `json:"avg"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" gorethink:"infected"`
	Result   string `json:"result" gorethink:"result"`
	Engine   string `json:"engine" gorethink:"engine"`
	Database string `json:"database" gorethink:"database"`
	Updated  string `json:"updated" gorethink:"updated"`
}

// ParseAVGOutput convert avg output into ResultsData struct
func ParseAVGOutput(avgout string, path string) (ResultsData, error) {

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
		log.Error("[ERROR] colonSeparated was empty: ", colonSeparated)
		log.Errorf("[ERROR] AVG output was: \n%s", avgout)
		// fmt.Println("[ERROR] colonSeparated was empty: ", colonSeparated)
		// fmt.Printf("[ERROR] AVG output was: \n%s", avgout)
		return ResultsData{}, errors.New("Unable to parse AVG output.")
	}

	return avg, nil
}

// Get Anti-Virus scanner version
func getAvgVersion() string {
	versionOut := utils.RunCommand("/usr/bin/avgscan", "-v")
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

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	utils.Assert(err)
	return string(updated)
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func updateAV() error {
	fmt.Println("Updating AVG...")
	// AVG needs to have the daemon started first
	exec.Command("/etc/init.d/avgd", "start").Output()

	fmt.Println(utils.RunCommand("avgupdate"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
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
	var elasitcsearch string
	var table bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:        "elasitcsearch",
			Value:       "",
			Usage:       "elasitcsearch address for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH",
			Destination: &elasitcsearch,
		},
		cli.BoolFlag{
			Name:        "table, t",
			Usage:       "output as Markdown table",
			Destination: &table,
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
			Action: func(c *cli.Context) error {
				return updateAV()
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			utils.Assert(err)
		}

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		// AVG needs to have the daemon started first
		exec.Command("/etc/init.d/avgd", "start").Output()
		// Give avgd a few to finish
		time.Sleep(time.Second * 2)

		var results ResultsData

		results, err := ParseAVGOutput(utils.RunCommand("/usr/bin/avgscan", path), path)
		if err != nil {
			// If fails try a second time
			results, err = ParseAVGOutput(utils.RunCommand("/usr/bin/avgscan", path), path)
			utils.Assert(err)
		}

		avg := AVG{
			Results: results,
		}

		// upsert into Database
		elasticsearch.InitElasticSearch()
		elasticsearch.WritePluginResultsToDatabase(elasticsearch.PluginResults{
			ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
			Name:     name,
			Category: category,
			Data:     structs.Map(avg.Results),
		})

		if table {
			printMarkDownTable(avg)
		} else {
			avgJSON, err := json.Marshal(avg)
			utils.Assert(err)
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
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
