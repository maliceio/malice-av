package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
	"github.com/crackcomm/go-clitable"
	"github.com/parnurzeal/gorequest"
)

// Version stores the plugin's version
var Version string

// BuildTime stores the plugin's build time
var BuildTime string

// Avira json object
type Avira struct {
	Results ResultsData `json:"avira"`
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

// ParseAviraOutput convert Avira output into ResultsData struct
func ParseAviraOutput(aviraout string) ResultsData {

	avira := ResultsData{Infected: false}

	fmt.Println(aviraout)

	return avira
}

// // extractVirusName extracts Virus name from scan results string
// func extractVirusName(line string) string {
// 	r := regexp.MustCompile(`<(.+)>`)
// 	res := r.FindStringSubmatch(line)
// 	if len(res) != 2 {
// 		return ""
// 	}
// 	return res[1]
// }

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

func printMarkDownTable(avira Avira) {

	fmt.Println("#### Avira")
	table := clitable.New([]string{"Infected", "Result", "Engine", "Updated"})
	table.AddRow(map[string]interface{}{
		"Infected": avira.Results.Infected,
		"Result":   avira.Results.Result,
		"Engine":   avira.Results.Engine,
		"Updated":  avira.Results.Updated,
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
	app.Name = "avira"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice Avira AntiVirus Plugin"
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
	app.Action = func(c *cli.Context) {
		path := c.Args().First()

		if _, err := os.Stat(path); os.IsNotExist(err) {
			assert(err)
		}
		cmdOut, err := exec.Command("/usr/lib/AntiVir/guard/avguard", "restart > /dev/null 2>&1").Output()
		if len(cmdOut) == 0 {
			assert(err)
		}

		avira := Avira{
			Results: ParseAviraOutput(RunCommand(
				"/usr/lib/AntiVir/guard/avscan",
				"-s",
				"--scan-in-archive=yes",
				"--scan-mode=all",
				"--heur-level=3",
				"--alert-action=none",
				"--heur-macro=yes",
				"--alert-action=none",
				path,
			)),
		}

		if c.Bool("table") {
			printMarkDownTable(avira)
		} else {
			aviraJSON, err := json.Marshal(avira)
			assert(err)
			if c.Bool("post") {
				request := gorequest.New()
				if c.Bool("proxy") {
					request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
				}
				request.Post(os.Getenv("MALICE_ENDPOINT")).
					Set("Task", path).
					Send(aviraJSON).
					End(printStatus)
			}
			fmt.Println(string(aviraJSON))
		}
	}

	err := app.Run(os.Args)
	assert(err)
}
