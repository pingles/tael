package tael

import (
	"encoding/json"
	"fmt"
	"text/template"
	"os"
	"strings"
	"github.com/fatih/color"
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

func NewTemplatedFormatter(layout string) *TemplatedFormatter {
	tmpl := template.Must(template.New("entry").Parse(layout))
	return &TemplatedFormatter{
		template: tmpl,
	}
}

func (f *TemplatedFormatter) Write(e *LogEntry) {
	err := f.template.Execute(os.Stdout, e)
	fmt.Println()
	if err != nil {
		fmt.Println("error templating entry:", err.Error())
		os.Exit(1)
	}
}
