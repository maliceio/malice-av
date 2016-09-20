package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "sophos"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"sophos" gorethink:"sophos"`
}

// Sophos json object
type Sophos struct {
	Results ResultsData `json:"sophos"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" gorethink:"infected"`
	Result   string `json:"result" gorethink:"result"`
	Engine   string `json:"engine" gorethink:"engine"`
	Database string `json:"database" gorethink:"database"`
	Updated  string `json:"updated" gorethink:"updated"`
}

// ParseSophosOutput convert sophos output into ResultsData struct
func ParseSophosOutput(sophosout string, path string) (ResultsData, error) {

	// root@0e01fb905ffb:/malware# savscan -f EICAR
	// SAVScan virus detection utility
	// Version 5.21.0 [Linux/AMD64]
	// Virus data version 5.27, April 2016
	// Includes detection for 11283995 viruses, Trojans and worms
	// Copyright (c) 1989-2016 Sophos Limited. All rights reserved.
	//
	// System time 03:48:15, System date 22 August 2016
	// Command line qualifiers are: -f
	//
	// Full Scanning
	//
	// >>> Virus 'EICAR-AV-Test' found in file EICAR
	//
	// 1 file scanned in 4 seconds.
	// 1 virus was discovered.
	// 1 file out of 1 was infected.
	// If you need further advice regarding any detections please visit our
	// Threat Center at: http://www.sophos.com/en-us/threat-center.aspx
	// End of Scan.

	version, database := getSophosVersion()

	sophos := ResultsData{
		Infected: false,
		Engine:   version,
		Database: database,
		Updated:  getUpdatedDate(),
	}

	lines := strings.Split(sophosout, "\n")

	for _, line := range lines {
		if strings.Contains(line, ">>> Virus") && strings.Contains(line, "found in file") {
			parts := strings.Split(line, "'")
			sophos.Result = strings.TrimSpace(parts[1])
			sophos.Infected = true
		}
	}

	return sophos, nil
}

// Get Anti-Virus scanner version
func getSophosVersion() (version string, database string) {
	// root@0e01fb905ffb:/malware# /opt/sophos/bin/savscan --version
	// SAVScan virus detection utility
	// Copyright (c) 1989-2016 Sophos Limited. All rights reserved.
	//
	// System time 03:41:05, System date 22 August 2016
	//
	// Product version           : 5.21.0
	// Engine version            : 3.64.0
	// Virus data version        : 5.27
	// User interface version    : 2.03.064
	// Platform                  : Linux/AMD64
	// Released                  : 26 April 2016
	// Total viruses (with IDEs) : 11283995
	versionOut := utils.RunCommand("/opt/sophos/bin/savscan", "--version")
	return parseSophosVersion(versionOut)
}

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
	fmt.Println("Updating Sophos...")
	// root@0e01fb905ffb:/opt/sophos/update# ./savupdate.sh
	// Updating from versions - SAV: 9.12.1, Engine: 3.64.0, Data: 5.27
	// Updating Sophos Anti-Virus....
	// Updating Command-line programs
	// Updating SAVScan on-demand scanner
	// Updating sav-protect startup script
	// Updating sav-rms startup script
	// Updating Virus Engine and Data
	// Updating Sophos Anti-Virus Scanning Daemon
	// Updating Talpa Kernel Support
	// Updating Manifest
	// Selecting appropriate kernel support...
	// On-access scanning not available because of problems during kernel support compilation.
	// Update completed.
	// Updated to versions - SAV: 9.12.2, Engine: 3.65.2, Data: 5.30
	// Successfully updated Sophos Anti-Virus from sdds:SOPHOS

	fmt.Println(utils.RunCommand("/opt/sophos/update/savupdate.sh"))

	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func printMarkDownTable(sophos Sophos) {

	fmt.Println("#### Sophos")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": sophos.Results.Infected,
		"Result":   sophos.Results.Result,
		"Engine":   sophos.Results.Engine,
		"Updated":  sophos.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

// writeToDatabase upserts plugin results into Database
func writeToDatabase(results pluginResults) {
	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", utils.Getopt("MALICE_RETHINKDB", "rethink")),
		Timeout:  5 * time.Second,
		Database: "malice",
	})
	if err != nil {
		log.Debug(err)
		return
	}
	defer session.Close()

	res, err := r.Table("samples").Get(results.ID).Run(session)
	utils.Assert(err)
	defer res.Close()

	if res.IsNil() {
		// upsert into RethinkDB
		resp, err := r.Table("samples").Insert(results, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
		utils.Assert(err)
		log.Debug(resp)
	} else {
		resp, err := r.Table("samples").Get(results.ID).Update(map[string]interface{}{
			"plugins": map[string]interface{}{
				category: map[string]interface{}{
					name: results.Data,
				},
			},
		}).RunWrite(session)
		utils.Assert(err)

		log.Debug(resp)
	}
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
	app.Name = "sophos"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Sophos AntiVirus Plugin"
	var rethinkdb string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:        "rethinkdb",
			Value:       "",
			Usage:       "rethinkdb address for Malice to store results",
			EnvVar:      "MALICE_RETHINKDB",
			Destination: &rethinkdb,
		},
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
		} else {
			r.Log.Out = ioutil.Discard
		}

		var results ResultsData

		results, err := ParseSophosOutput(utils.RunCommand("savscan", "-f", path), path)
		if err != nil {
			// If fails try a second time
			results, err = ParseSophosOutput(utils.RunCommand("savscan", "-f", path), path)
			utils.Assert(err)
		}

		// upsert into Database
		writeToDatabase(pluginResults{
			ID:   utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
			Data: results,
		})

		sophos := Sophos{
			Results: results,
		}

		if c.Bool("table") {
			printMarkDownTable(sophos)
		} else {
			sophosJSON, err := json.Marshal(sophos)
			utils.Assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(sophosJSON).
					End(printStatus)
			}
			fmt.Println(string(sophosJSON))
		}
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
