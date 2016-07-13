package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/addons/search"
	"github.com/inpime/fader/store"
)

func AppStoresInitFromExistBuckets() error {
	for _, bucket := range search.GetAllBuckets() {

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
