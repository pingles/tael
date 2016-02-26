package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/fatih/color"
	"os"
	"strings"
	"text/template"
)

var (
	host            = kingpin.Flag("host", "aws elasticsearch url").OverrideDefaultFromEnvar("ELASTICSEARCH_HOST").String()
	follow          = kingpin.Flag("follow", "follow log").Short('f').Bool()
	index           = kingpin.Flag("index", "elasticsearch index").Default("*").String()
	numberOfResults = kingpin.Flag("number", "number of results to retrieve").Default("10").Short('n').Int()
	outputJson      = kingpin.Flag("json", "output as json").Short('j').Bool()
	layout          = kingpin.Flag("layout", "custom templated output").Short('l').String()
	query           = kingpin.Arg("query", "elasticsearch query").String()
)

type Formatter interface {
	Write(e *LogEntry)
}

type JsonFormatter struct{}

func (j *JsonFormatter) Write(e *LogEntry) {
	bytes, err := json.Marshal(e)
	if err != nil {
		fmt.Println("error marshaling to json:", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(bytes))
}

type StandardFormatter struct{}

func colourLevel(e *LogEntry) string {
	lname := strings.ToUpper(e.LevelName)
	if strings.Contains(lname, "INFO") {
		return color.GreenString("%-8s", lname)
	}
	if strings.Contains(lname, "WARN") {
		return color.YellowString("%-8s", lname)
	}
	if strings.Contains(lname, "ERROR") {
		return color.RedString("%-8s", lname)
	}
	if strings.Contains(lname, "DEBUG") {
		return color.BlueString("%-8s", lname)
	}

	return fmt.Sprintf("%-8s", e.LevelName)
}

func (f *StandardFormatter) Write(e *LogEntry) {
	time := color.CyanString(e.Time.Format("2006-01-02 15:04:05.000"))
	line := fmt.Sprintf("%s: %23s %s %s", color.BlueString(e.Id), time, colourLevel(e), e.Message)
	fmt.Println(line)
}

type TemplatedFormatter struct {
	template *template.Template
}

func (f *TemplatedFormatter) Write(e *LogEntry) {
	err := f.template.Execute(os.Stdout, e)
	fmt.Println()
	if err != nil {
		fmt.Println("error templating entry:", err.Error())
		os.Exit(1)
	}
}

func createFormatter() Formatter {
	if *outputJson {
		return &JsonFormatter{}
	} else if *layout != "" {
		tmpl := template.Must(template.New("entry").Parse(*layout))
		return &TemplatedFormatter{
			template: tmpl,
		}
	}
	return &StandardFormatter{}
}

func main() {
	kingpin.Parse()
	if *query == "" {
		kingpin.FatalUsage("query cannot be blank.")
	}
	if *host == "" {
		kingpin.FatalUsage("host cannot be blank.")
	}

	entries := PerformSearch(*index, *query, *host)

	n := 0

	formatter := createFormatter()
	for e := range entries {
		formatter.Write(e)
		n += 1

		if !*follow && n == *numberOfResults {
			return
		}
	}
}
