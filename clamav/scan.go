package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	name     = "clamav"
	category = "av"
)

type pluginResults struct {
	ID   string      `json:"id" gorethink:"id,omitempty"`
	Data ResultsData `json:"clamav" gorethink:"clamav"`
}

// ClamAV json object
type ClamAV struct {
	Results ResultsData `json:"clamav"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" gorethink:"infected"`
	Result   string `json:"result" gorethink:"result"`
	Engine   string `json:"engine" gorethink:"engine"`
	Known    string `json:"known" gorethink:"known"`
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

// ParseClamAvOutput convert clamav output into ClamAV struct
func ParseClamAvOutput(clamout string) ResultsData {

	clamAV := ResultsData{}

	lines := strings.Split(clamout, "\n")
	// Extract AV Scan Result
	result := lines[0]
	if len(result) != 0 {
		pathAndResult := strings.Split(result, ":")
		if strings.Contains(pathAndResult[1], "OK") {
			clamAV.Infected = false
		} else {
			clamAV.Infected = true
			clamAV.Result = strings.TrimSpace(strings.TrimRight(pathAndResult[1], "FOUND"))
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
					clamAV.Known = strings.TrimSpace(keyvalue[1])
				case strings.Contains(line, "Engine version"):
					clamAV.Engine = strings.TrimSpace(keyvalue[1])
				}
			}
		}
	}

	clamAV.Updated = getUpdatedDate()

	return clamAV
}

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return string(updated)
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func printMarkDownTable(clamav ClamAV) {
	fmt.Println("#### ClamAV")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": clamav.Results.Infected,
		"Result":   clamav.Results.Result,
		"Engine":   clamav.Results.Engine,
		// "Known":    clamav.Results.Known,
		"Updated": clamav.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

func updateAV() error {
	fmt.Println("Updating ClamAV...")
	fmt.Println(RunCommand("freshclam"))
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
	app.Name = "clamav"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice ClamAV Plugin"
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
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		clamav := ClamAV{
			Results: ParseClamAvOutput(RunCommand("/usr/bin/clamscan", "--stdout", path)),
		}

		// upsert into Database
		writeToDatabase(pluginResults{ID: getSHA256(path), Data: clamav.Results})

		if c.Bool("table") {
			printMarkDownTable(clamav)
		} else {
			fprotJSON, err := json.Marshal(clamav)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(fprotJSON).
					End(printStatus)
			}
			fmt.Println(string(fprotJSON))
		}
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
