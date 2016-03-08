package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/pingles/tael"
	"time"
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

func createFormatter() tael.Formatter {
	if *outputJson {
		return &tael.JsonFormatter{}
	} else if *layout != "" {
		return tael.NewTemplatedFormatter(*layout)
	}
	return &tael.StandardFormatter{}
}

func main() {
	kingpin.Parse()
	if *query == "" {
		kingpin.FatalUsage("query cannot be blank.")
	}
	if *host == "" {
		kingpin.FatalUsage("host cannot be blank.")
	}

	search := tael.NewSearch(*index, *query, time.Now())
	entries := tael.PerformSearch(*host, search)

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
