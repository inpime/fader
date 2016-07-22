package store

import (
	"github.com/inpime/fader/utils/sdata"
	"gopkg.in/olivere/elastic.v3"
)

type SearchQueryIface interface {
	QueryString() string
	SetQueryString(string) interface{}

	Page() int
	SetPage(int) interface{}

	PerPage() int
	SetPerPage(int) interface{}

	Bucket() string
	SetBucket(string) interface{}

	QueryRaw() map[string]interface{}
	SetQueryRaw(map[string]interface{}) interface{}
}

type SearchResultIface interface {
	GetFiles() []*File
	GetAggs() *elastic.Aggregations

	GetTotal() int64
	GetHasNext() bool
	GetNextPage() int
	IsError() bool
	Error() string
}

type Searcher interface {
	UpdateDocument(id string, v interface{}) error
	DeleteDocument(id string) error
	UpdateMapping(bucketName string, mapping map[string]interface{}) error
	Search(SearchQueryIface) SearchResultIface
}

type FileSearchDTO struct {
	FileID   string
	Bucket   string
	Name     string
	Meta     map[string]interface{}
	MapData  map[string]interface{}
	TextData string
	// CreatedAt time.Time
	// UpdatedAt time.Time
}

func FileSearchFromFile(file *File) *FileSearchDTO {
	dto := &FileSearchDTO{
		FileID:  file.ID(),
		Name:    file.Name(),
		Bucket:  file.Bucket(),
		MapData: file.MapData(),
		Meta:    file.Meta(),
		// CreatedAt: file.CreatedAt(),
		// UpdatedAt: file.UpdatedAt(),
	}

	if file.IsText() {
		dto.TextData = string(file.RawData().Bytes())
	}

	return dto
}

func UpdateSearchMapping(bucket *Bucket) error {
	type fieldMapping struct {
		Type   string `json:"type"`
		Index  string `json:"index,omitempty"`
		Format string `json:"format,omitempty"`
	}

	docMapping := sdata.NewStringMapFrom(BucketFileMappingDefault).
		Set("Meta", map[string]interface{}{"properties": bucket.MetaDataFilesMapping()})

	if len(bucket.MapDataFilesMapping()) > 0 {

		docMapping.Set("MapData", map[string]interface{}{"properties": bucket.MapDataFilesMapping()})
	}

	mapping := map[string]interface{}{
		"properties": docMapping.ToMap(),
	}

	// b, errr := json.Marshal(mapping)
	// logrus.Infof("Mapping %q, %v, \n%s\n", bucket.Name(), errr, string(b))

	_, err := ESClient.
		PutMapping().
		Index(ElasticSearchIndexName).
		Type(bucket.Name()).
		BodyJson(mapping).
		Do()

	return err
}

func UpdateSearchDocument(typeName, id string, v interface{}) error {
	_, err := ESClient.Index().
		Index(ElasticSearchIndexName).
		Type(typeName).
		Id(id).
		BodyJson(v).
		Do()

	return err
}

func DeleteSearchDocument(typeName, id string) error {
	_, err := ESClient.Delete().
		Index(ElasticSearchIndexName).
		Type(typeName).
		Id(id).
		Do()

	return err
}
