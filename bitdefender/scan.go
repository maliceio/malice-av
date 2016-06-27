package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/crackcomm/go-clitable"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
	r "gopkg.in/dancannon/gorethink.v2"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

const (
	name     = "bitdefender"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"bitdefender" gorethink:"bitdefender"`
}

// Bitdefender json object
type Bitdefender struct {
	Results ResultsData `json:"bitdefender"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" gorethink:"infected"`
	Result   string `json:"result" gorethink:"result"`
	Engine   string `json:"engine" gorethink:"engine"`
	Updated  string `json:"updated" gorethink:"updated"`
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

// getSHA256 calculates a file's sha256sum
func getSHA256(name string) string {

	dat, err := ioutil.ReadFile(name)
	assert(err)

	h256 := sha256.New()
	_, err = h256.Write(dat)
	assert(err)

	return fmt.Sprintf("%x", h256.Sum(nil))
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

	bitdefender := ResultsData{Infected: false}
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

	bitdefender.Updated = getUpdatedDate()

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

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return string(updated)
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

func updateAV() error {
	fmt.Println("Updating Bitdefender...")
	fmt.Println(RunCommand("bdscan", "--update"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

// writeToDatabase upserts plugin results into Database
func writeToDatabase(results pluginResults) {

	address := fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink"))

	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  address,
		Timeout:  5 * time.Second,
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Get(results.ID).Run(session)
		assert(err)
		defer res.Close()

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(results, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			assert(err)
			log.Debug(resp)
		} else {
			resp, err := r.Table("samples").Get(results.ID).Update(map[string]interface{}{
				"plugins": map[string]interface{}{
					category: map[string]interface{}{
						name: results.Data,
					},
				},
			}).RunWrite(session)
			assert(err)

			log.Debug(resp)
		}

	} else {
		log.Debug(err)
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
	app.Name = "bitdefender"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Bitdefender AntiVirus Plugin"
	var rethinkdb string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
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
		cli.StringFlag{
			Name:        "rethinkdb",
			Value:       "",
			Usage:       "rethinkdb address for Malice to store results",
			EnvVar:      "MALICE_RETHINKDB",
			Destination: &rethinkdb,
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
		path, err := filepath.Abs(c.Args().First())
		assert(err)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}
		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		} else {
			r.Log.Out = ioutil.Discard
		}

		bitdefender := Bitdefender{
			Results: ParseBitdefenderOutput(RunCommand("bdscan", path)),
		}

		// upsert into Database
		writeToDatabase(pluginResults{
			ID:   getopt("MALICE_SCANID", getSHA256(path)),
			Data: bitdefender.Results,
		})

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
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
