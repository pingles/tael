package elasticsearch

import (
	"fmt"
	"path"
	"io"
	"encoding/json"
	"bytes"
)


// mirrors the elasticsearch search json structures

type MustFilter interface {
	MarshalJSON() ([]byte, error)
}

const (
	RangeFormat = "epoch_millis"
)

type MustRange struct {
	Field  string
	Format string
	Lte    interface{}
	Gte    interface{}
}

func (q *MustRange) MarshalJSON() ([]byte, error) {
	m := map[string]interface{} {
		"range": map[string]interface{} {
			q.Field: map[string]interface{} {
				"format": q.Format,
				"lte": q.Lte,
				"gte": q.Gte,
			},
		},
	}
	return json.Marshal(m)
}


const (
	MatchPhrase = "phrase"
)

type MustQuery struct {
	Field     string
	MatchType string
	Query     string
}

func (q *MustQuery) MarshalJSON() ([]byte, error) {
	m := map[string]interface{} {
		"query": map[string]interface{} {
			"match": map[string]interface{} {
				q.Field: map[string]interface{} {
					"type": q.MatchType,
					"query": q.Query,
				},
			},
		},
	}
	return json.Marshal(m)
}


type Bool struct {
	Must []MustFilter `json:"must"`
}

type Filter struct {
	Bool *Bool `json:"bool"`
}

type FilterQuery struct {
	QueryString string
}

func (q *FilterQuery) MarshalJSON() ([]byte, error) {
	m := map[string]interface{} {
		"query_string": map[string]interface{} {
			"analyze_wildcard": true,
			"query": q.QueryString,
		},
	}
	return json.Marshal(m)
}

type FilteredQuery struct {
	Filter  *Filter      `json:"filter"`
	Query   *FilterQuery `json:"query"`
}

type Query struct {
	Filtered *FilteredQuery `json:"filtered"`
}

const (
	Descending = "desc"
	Ascending = "asc"
)

type Sort struct {
	Field  string
	Order  string
}

func (s *Sort) MarshalJSON() ([]byte, error) {
	m := map[string]interface{} {
		s.Field: map[string]interface{} {
			"unmapped_type": "boolean",
			"order": s.Order,
		},
	}
	return json.Marshal(m)
}


type Search struct {
	Query *Query `json:"query"`
	Sort  *Sort  `json:"sort"`
	Size  int32  `json:"size"`
}



// TODO
// think of a better name for this
type SearchContext struct {
	Host        string
	Index       string

	// query will be parsed using Lucene syntax
	Search      *Search
}

func (s *SearchContext) url() string {
	return fmt.Sprintf("%s/%s", s.Host, path.Join(s.Index, "_search"))
}

func (s *SearchContext) body() (io.Reader, error) {
	bs, err := json.Marshal(s.Search)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bs), nil
}

const (
	DefaultSize = 20
)

func NewSearchContext(host, index string, search *Search) *SearchContext {
	return &SearchContext{
		Host: host,
		Index: index,
		Search: search,
	}
}
