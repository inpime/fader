package api

import (
	"github.com/Sirupsen/logrus"
	"store"
)

func makeSearch(filter store.SearchQueryIface) store.SearchResultIface {
	result := &store.SearchResult{}

	bucket, err := store.BucketByName(filter.Bucket())

	if err != nil {
		result.Err = err
		return result
	}

	query := store.MustSearchService()
	query.Index(Cfg.Search.IndexName)
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
					"file_id": hit.Id,
					"bucket":  bucket.Name(),
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
