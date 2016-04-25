package main

import (
	"errors"
	"gopkg.in/alecthomas/kingpin.v2"
	es "github.com/pingles/tael/elasticsearch"
	"time"
)

var (
	host            = kingpin.Flag("host", "aws elasticsearch url").OverrideDefaultFromEnvar("ELASTICSEARCH_HOST").String()
	follow          = kingpin.Flag("follow", "follow log").Short('f').Bool()
	index           = kingpin.Flag("index", "elasticsearch index").Default("*").String()
	numberOfResults = kingpin.Flag("number", "number of results to retrieve").Default("10").Short('n').Int()
	outputJson      = kingpin.Flag("json", "output as json").Short('j').Bool()
	layout          = kingpin.Flag("layout", "custom templated output").Short('l').String()
	query           = kingpin.Flag("query", "elasticsearch query").Default("*").String()
	filters         = kingpin.Arg("filter", "field filter, name=value").StringMap()
)

func createFormatter() Formatter {
	if *outputJson {
		return &JsonFormatter{}
	} else if *layout != "" {
		return NewTemplatedFormatter(*layout)
	}
	return &StandardFormatter{}
}

type logMessage struct {
	Message   string
	Timestamp time.Time
	LevelName string
	LogName   string
}
func newLogMessageFromHit(hit *es.Hit) (*logMessage, error) {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", hit.Source["@timestamp"].(string))
	if err != nil {
		return nil, err
	}

	message, ok := hit.Source["message"]
	if !ok {
		return nil, errors.New("no message")
	}
	log, ok := hit.Source["logname"]
	if !ok {
		return nil, errors.New("no logname")
	}
	level, ok := hit.Source["level_name"]
	if !ok {
		return nil, errors.New("no level")
	}

	return &logMessage{
		Message: message.(string),
		Timestamp: t,
		LogName: log.(string),
		LevelName: level.(string),
	}, nil
}

func main() {
	kingpin.Parse()
	if *host == "" {
		kingpin.FatalUsage("host cannot be blank.")
	}

	mustFilters := make([]es.MustFilter, len(*filters))
	i := 0
	for fieldName, fieldQuery := range *filters {
		mustFilters[i] = &es.MustQuery{
			Field: fieldName,
			MatchType: es.MatchPhrase,
			Query: fieldQuery,
		}
		i = i + 1
	}

	search := &es.Search{
		Size: es.DefaultSize,
		Sort: &es.Sort{
			Field: "@timestamp",
			Order: es.Descending,
		},
		Query: &es.Query{
			Filtered: &es.FilteredQuery{
				Filter: &es.Filter{
					Bool: &es.Bool{
						Must: mustFilters,
					},
				},
				Query: &es.FilterQuery{
					QueryString: *query,
				},
			},
		},
	}
	s := es.NewSearchContext(*host, *index, search)
	formatter := createFormatter()
	for hit := range es.StreamSearch(s) {
		msg, err := newLogMessageFromHit(hit)
		if err != nil {
			continue
		}
		formatter.Write(msg)
	}
}
