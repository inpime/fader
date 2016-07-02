package search

import (
	"api/config"
	"store"
)

func GetAllBuckets() []*store.File {
	// all buckets
	filter := store.NewSearchFilter(config.BucketsBucketName)
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(100)

	queryRaw := BuildSearchQueryFilesByBucket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return MakeSearch(filter).GetFiles()
}

func GetAllFiles(bucket string) []*store.File {
	// all buckets
	filter := store.NewSearchFilter(bucket)
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(1000)

	queryRaw := BuildSearchQueryFilesByBucket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return MakeSearch(filter).GetFiles()
}
