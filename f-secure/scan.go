package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	name     = "avast"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"avast" gorethink:"avast"`
}

// Avast json object
type Avast struct {
	Results ResultsData `json:"avast"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" gorethink:"infected"`
	Result   string `json:"result" gorethink:"result"`
	Engine   string `json:"engine" gorethink:"engine"`
	Database string `json:"database" gorethink:"database"`
	Updated  string `json:"updated" gorethink:"updated"`
}

// ParseAvastOutput convert avast output into ResultsData struct
func ParseAvastOutput(avastout string, path string) (ResultsData, error) {

	avast := ResultsData{
		Infected: false,
		Engine:   getAvastVersion(),
		Database: getAvastVPS(),
		Updated:  getUpdatedDate(),
	}

	result := strings.Split(avastout, "\t")

	if !strings.Contains(avastout, "[OK]") {
		avast.Infected = true
		avast.Result = strings.TrimSpace(result[1])
	}

	return avast, nil
}

// Get Anti-Virus scanner version
func getAvastVersion() string {
	versionOut := utils.RunCommand("/bin/scan", "-v")
	return strings.TrimSpace(versionOut)
}

func getAvastVPS() string {
	versionOut := utils.RunCommand("/bin/scan", "-V")
	return strings.TrimSpace(versionOut)
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
	fmt.Println("Updating Avast...")
	// Avast needs to have the daemon started first
	exec.Command("/etc/init.d/avast", "start").Output()

	fmt.Println(utils.RunCommand("/var/lib/avast/Setup/avast.vpsupdate"))
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err := ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func printMarkDownTable(avast Avast) {

	fmt.Println("#### Avast")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": avast.Results.Infected,
		"Result":   avast.Results.Result,
		"Engine":   avast.Results.Engine,
		"Updated":  avast.Results.Updated,
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
	app.Name = "avast"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Avast AntiVirus Plugin"
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

		// Avast needs to have the daemon started first
		exec.Command("/etc/init.d/avast", "start").Output()
		// Give avast service a few to finish
		// time.Sleep(time.Second * 2)

		var results ResultsData

		results, err := ParseAvastOutput(utils.RunCommand("fsav", "--virus-action=none", path), path)
		if err != nil {
			// If fails try a second time
			results, err = ParseAvastOutput(utils.RunCommand("fsav", "--virus-action=none", path), path)
			utils.Assert(err)
		}

		// upsert into Database
		// database.WriteToDatabase(pluginResults{
		writeToDatabase(pluginResults{
			ID:   utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
			Data: results,
		})

		avast := Avast{
			Results: results,
		}

		if c.Bool("table") {
			printMarkDownTable(avast)
		} else {
			avastJSON, err := json.Marshal(avast)
			utils.Assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(avastJSON).
					End(printStatus)
			}
			fmt.Println(string(avastJSON))
		}
		return nil
	}

	err := app.Run(os.Args)
	utils.Assert(err)
}
