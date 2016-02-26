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
	"github.com/levigross/grequests"
	"github.com/parnurzeal/gorequest"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// Comodo json object
type Comodo struct {
	Results ResultsData `json:"comodo"`
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

// ParseComodoOutput convert comodo output into ResultsData struct
func ParseComodoOutput(comodoout string) ResultsData {

	comodo := ResultsData{Infected: false, Engine: "1.1", Updated: BuildTime}
	// -----== Scan Start ==-----
	// /malware/EICAR ---> Found Virus, Malware Name is Malware
	// -----== Scan End ==-----
	// Number of Scanned Files: 1
	// Number of Found Viruses: 1

	lines := strings.Split(comodoout, "\n")
	fmt.Println(lines)
	// Extract Virus string and extract colon separated lines into an slice
	for _, line := range lines {
		if len(line) != 0 {
			if strings.Contains(line, "Found Virus") {
				result := extractVirusName(line)
				if len(result) != 0 {
					comodo.Result = result
					comodo.Infected = true
				} else {
					fmt.Println("[ERROR] Virus name extracted was empty: ", result)
					os.Exit(2)
				}
			}
		}
	}

	return comodo
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	keyvalue := strings.Split(line, "is")
	return strings.TrimSpace(keyvalue[1])
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func parseUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func printMarkDownTable(comodo Comodo) {

	fmt.Println("#### Comodo")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": comodo.Results.Infected,
		"Result":   comodo.Results.Result,
		"Engine":   comodo.Results.Engine,
		"Updated":  comodo.Results.Updated,
	})
	table.Markdown = true
	table.Print()
}

func updateAV() {
	fmt.Println("Updating Comodo...")
	response, err := grequests.Get("http://download.comodo.com/av/updates58/sigs/bases/bases.cav", nil)
	assert(err)

	if response.Ok != true {
		log.Println("Request did not return OK")
	}

	if err := response.DownloadToFile("/opt/COMODO/scanners/bases.cav"); err != nil {
		log.Println("Unable to download file: ", err)
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
	app.Name = "comodo"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Comodo AntiVirus Plugin"
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

		comodo := Comodo{
			Results: ParseComodoOutput(RunCommand("/opt/COMODO/cmdscan", "-vs", path)),
		}

		if c.Bool("table") {
			printMarkDownTable(comodo)
		} else {
			comodoJSON, err := json.Marshal(comodo)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(comodoJSON).
					End(printStatus)
			}
			fmt.Println(string(comodoJSON))
		}
	}

	err := app.Run(os.Args)
	assert(err)
}
