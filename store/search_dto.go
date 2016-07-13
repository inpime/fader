package store

import (
	"gopkg.in/olivere/elastic.v3"
)

var _ SearchResultIface = (*SearchResult)(nil)
var _ SearchQueryIface = (*SearchQuery)(nil)

// ------------------
// Search query
// ------------------

func NewSearchFilter(bucketName string) *SearchQuery {
	return &SearchQuery{
		QueryStr:   "",
		PageNum:    0,
		PerPageNum: 25,
		BucketName: bucketName,

		Query: make(map[string]interface{}),
	}
}

type SearchQuery struct {
	QueryStr   string
	PageNum    int
	PerPageNum int
	BucketName string

	Query map[string]interface{}
}

func (s SearchQuery) QueryString() string {
	return s.QueryStr
}

func (s *SearchQuery) SetQueryString(str string) interface{} {
	s.QueryStr = str
	return s
}

func (s SearchQuery) Page() int {
	return s.PageNum
}

func (s *SearchQuery) SetPage(i int) interface{} {
	s.PageNum = i
	return s
}

func (s SearchQuery) PerPage() int {
	return s.PerPageNum
}

func (s *SearchQuery) SetPerPage(i int) interface{} {
	s.PerPageNum = i
	return s
}

func (s SearchQuery) Bucket() string {
	return s.BucketName
}

func (s *SearchQuery) SetBucket(str string) interface{} {
	s.BucketName = str
	return s
}

func (s SearchQuery) QueryRaw() map[string]interface{} {
	return s.Query
}

func (s *SearchQuery) SetQueryRaw(v map[string]interface{}) interface{} {
	s.Query = v
	return s
}

// ------------------
// Search result
// ------------------

type SearchResult struct {
	Files []*File
	Aggs  *elastic.Aggregations `,omitempty"`

	Total       int64
	HasNext     bool
	NextPage    int
	CurrentPage int
	PerPage     int

	Err error
}

func (s *SearchResult) GetFiles() []*File {
	return s.Files
}

func (s SearchResult) GetAggs() *elastic.Aggregations {
	return s.Aggs
}

func (s SearchResult) GetTotal() int64 {
	return s.Total
}

func (s SearchResult) GetHasNext() bool {
	return s.HasNext
}

func (s SearchResult) GetNextPage() int {
	return s.NextPage
}

func (s SearchResult) GetCurrentPage() int {
	return s.CurrentPage
}

func (s SearchResult) IsError() bool {
	return s.Err != nil
}

func (s SearchResult) Error() string {
	return s.Err.Error()
}
