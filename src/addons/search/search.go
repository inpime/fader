package search

import (
	"github.com/Sirupsen/logrus"
	"store"
	"strings"
	"utils"
)

func makeSearch(filter store.SearchQueryIface) store.SearchResultIface {
	result := &store.SearchResult{}

	bucket, err := store.BucketByName(filter.Bucket())

	if err != nil {
		result.Err = err
		return result
	}

	query := store.MustSearchService()
	query.Index("fader")
	query.Type(bucket.Name())
	query.SetQueryRaw(filter.QueryRaw())

	searchResult, err := query.Do()

	if err != nil {
		result.Err = err
		return result
	}

	result.Total = searchResult.Hits.TotalHits
	result.Aggs = &searchResult.Aggregations
	result.CurrentPage = filter.Page()
	result.PerPage = filter.PerPage()

	if searchResult.Hits != nil {
		for _, hit := range searchResult.Hits.Hits {
			file, err := store.LoadOrNewFileIDViaStores(
				bucket.Name(),
				hit.Id,
				bucket.MetaDataStoreName(),
				bucket.MapDataStoreName(),
				bucket.RawDataStoreName())

			if err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"_service": addonName,
					"file_id":  hit.Id,
					"bucket":   bucket.Name(),
				}).Warning("load file")
			}

			result.Files = append(result.Files, file)
		}
	}

	if searchResult.Hits != nil && searchResult.Hits.TotalHits > int64(filter.PerPage()*filter.Page()+len(searchResult.Hits.Hits)) {

		result.NextPage = filter.Page() + 1
		result.HasNext = true
	}

	return result
}

// buildSearchQueryFilesByBycket
// поиск всех файлов в указанном бакете по текстовому запросу (по умолчанию по RawData)
// Sorted by CreatedAt desc
func buildSearchQueryFilesByBycket(bucketName, qstr string, page, perpage int) map[string]interface{} {
	query := utils.Map()

	// filter arguments

	qstr = strings.TrimSpace(qstr)
	bucketName = strings.TrimSpace(bucketName)

	from := perpage * page
	if page <= 0 {
		from = 0
	}

	size := perpage

	// prepare arguments of query

	queryFileter := utils.Map().
		Set("term", utils.Map().Set("Bucket", bucketName))

	querySort := utils.Map().
		Set("Meta.CreatedAt", utils.Map().Set("order", "desc"))

	_query := utils.Map().
		Set("query_string", utils.Map().
			Set("default_field", "TextData").
			Set("query", qstr))

	// build query

	// 1. sort
	// 2. from
	// 3. size
	// 4. filter
	// 5. fields
	// 6. query

	// 1. sort
	query.Set("sort", []interface{}{querySort})

	// 2. from
	query.Set("from", from)

	// 3. size
	query.Set("size", size)

	// 4. filter
	query.Set("filter", queryFileter)

	// 5. filter
	query.Set("fields", []string{})

	// 6. filter
	if len(qstr) > 0 {
		query.Set("query", _query)
	}

	return query
}
