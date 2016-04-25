package main

import (
	"github.com/pingles/tael"
	es "github.com/pingles/tael/elasticsearch"
	"gopkg.in/alecthomas/kingpin.v2"
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

func main() {
	kingpin.Parse()
	if *host == "" {
		kingpin.FatalUsage("host cannot be blank.")
	}

	mustFilters := make([]es.MustFilter, len(*filters))
	i := 0
	for fieldName, fieldQuery := range *filters {
		mustFilters[i] = &es.MustQuery{
			Field:     fieldName,
			MatchType: es.MatchPhrase,
			Query:     fieldQuery,
		}
		i = i + 1
	}

	earliest := time.Now().Add(-15 * time.Minute)
	rangeFilter := &es.MustRange{
		Field:  "@timestamp",
		Format: es.RangeTimeEpoch,
		Gte:    es.EpochMillis(earliest),
	}
	mustFilters = append(mustFilters, rangeFilter)

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
		msg, err := tael.NewLogMessageFromHit(hit)
		if err != nil {
			continue
		}
		formatter.Write(msg)
	}
}
