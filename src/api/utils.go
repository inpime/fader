package api

import (
	"api/config"
	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"store"
)

func getAllBuckets() []*store.File {
	// all buckets
	filter := store.NewSearchFilter(config.BucketsBucketName)
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(100)

	queryRaw := buildSearchQueryFilesByBycket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return makeSearch(filter).GetFiles()
}

func getAllFiles(bucket string) []*store.File {
	// all buckets
	filter := store.NewSearchFilter(bucket)
	filter.SetQueryString("")
	filter.SetPage(0)
	filter.SetPerPage(1000)

	queryRaw := buildSearchQueryFilesByBycket(
		filter.Bucket(),
		filter.QueryString(),
		filter.Page(),
		filter.PerPage(),
	)
	filter.SetQueryRaw(queryRaw)

	return makeSearch(filter).GetFiles()
}

func AppStoresInitFromExistBuckets() error {
	for _, bucket := range getAllBuckets() {

		_bucket := &store.Bucket{&dbox.Bucket{*bucket.File}}

		// logrus.Debugf("settings: bucket %q, meta %#v", _bucket.Name(), _bucket.Meta())

		logrus.Debugf("settings: bucket %q, meta store type %#v, without postfix %v",
			_bucket.Name(),
			_bucket.GetMetaDataStoreType(),
			_bucket.GetMetaDataStoreNameWithoutPostfix(),
		)
		logrus.Debugf("settings: bucket %q, map store type %#v, without postfix %v",
			_bucket.Name(),
			_bucket.GetMapDataStoreType(),
			_bucket.GetMapDataStoreNameWithoutPostfix(),
		)
		logrus.Debugf("settings: bucket %q, raw store type %#v, without postfix %v",
			_bucket.Name(),
			_bucket.GetRawDataStoreType(),
			_bucket.GetRawDataStoreNameWithoutPostfix(),
		)

		_bucket.InitMetaDataStore(_bucket.GetMetaDataStoreType(),
			_bucket.GetMetaDataStoreNameWithoutPostfix())
		_bucket.InitMapDataStore(_bucket.GetMapDataStoreType(),
			_bucket.GetMapDataStoreNameWithoutPostfix())
		_bucket.InitRawDataStore(_bucket.GetRawDataStoreType(),
			_bucket.GetRawDataStoreNameWithoutPostfix())
	}

	return nil
}
