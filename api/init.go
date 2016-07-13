package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/addons/importexport"
	"github.com/inpime/fader/store"
	"gopkg.in/olivere/elastic.v3"
	"time"

	"github.com/inpime/fader/api/config"
	"net"
	"os"
)

func Init() error {
	if err := initElasticSearch(); err != nil {
		return err
	}

	db, err := bolt.Open(config.Cfg.Store.BoltDBFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		logrus.WithField("ref", "api").Error("Bolt db:", err)
		return err
	}

	// set default buckets store for buckets
	dbox.BucketStore = store.NewBoltDBStore(db, config.BucketsBucketName)
	store.BoltDBClient = db

	// root bucket buckets
	bucket, err := store.BucketByName(config.BucketsBucketName)
	bucket.InitInOneStore(dbox.BoltDBStoreType)

	isNewInstallation := false

	if err == dbox.ErrNotFound {
		logrus.WithField("ref", "api").
			Infof("buckets: create bucket %q\n", bucket.Name())

		bucket.MMetaDataFilesMapping().LoadFrom(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().LoadFrom(store.BucketMapMappingDefault)

		bucket.UpdateMapping()
		bucket.Sync()

		isNewInstallation = true
	}

	if isNewInstallation {
		logrus.WithField("ref", "api").
			Info("The first run. Installation of the console panel...")

		if err := importexport.AppImportFromLastArchive(importexport.ArchiveURLLatestVersion); err != nil {
			logrus.WithField("ref", "api").Error(err)
			return err
		}
	} else {
		logrus.Info("Existing application. Initializing settings....")
		if err := AppStoresInitFromExistBuckets(); err != nil {
			logrus.WithField("ref", "api").Error(err)
			return err
		}
	}

	return nil
}

func initElasticSearch() error {
	host := os.Getenv("FADER_ESADDRESSDOCKER")

	if len(host) > 0 {
	}

	if len(host) > 0 {
		logrus.WithField("ref", "api").
			Infof("Elasticsearch via docker. Host: %q", host)
		ips, err := net.LookupIP(host)
		if err != nil {
			logrus.WithField("ref", "api").
				Error(err)
			return err
		}
		logrus.WithField("ref", "api").
			Infof("Lookup for elasticsearch returns the following IPs:")
		for _, ip := range ips {
			config.Cfg.Search.Host = "http://" + ip.String() + ":9200"
			logrus.WithField("ref", "api").
				Infof("%v", ip)
			break
		}
	}

	var esLoggerOption elastic.ClientOptionFunc
	switch logrus.GetLevel() {
	case logrus.InfoLevel:
		esLoggerOption = elastic.SetInfoLog(logrus.StandardLogger().WithField("ref", "es"))
	case logrus.DebugLevel:
		esLoggerOption = elastic.SetTraceLog(logrus.StandardLogger().WithField("ref", "es"))
	case logrus.WarnLevel, logrus.ErrorLevel:
		esLoggerOption = elastic.SetErrorLog(logrus.StandardLogger().WithField("ref", "es"))
	}

	db, err := elastic.NewClient(
		// elastic.SetSniff(false),
		elastic.SetURL(config.Cfg.Search.Host),
		esLoggerOption,
		elastic.SetHealthcheckTimeoutStartup(time.Second*60),
		// elastic.SetTraceLog(logrus.New()),
	)

	if err != nil {
		logrus.WithField("ref", "api").
			Error("Setup ES:", err)
		return err
	}

	exists, err := db.IndexExists(store.ElasticSearchIndexName).Do()

	if err != nil {
		logrus.WithField("ref", "api").
			WithError(err).
			Errorf("exsist index %q", store.ElasticSearchIndexName)
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := db.
			CreateIndex(store.ElasticSearchIndexName).
			Do()

		if err != nil {
			// Handle error
			logrus.WithField("ref", "api").
				WithError(err).
				Errorf("create index %q", store.ElasticSearchIndexName)
			return err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged

			logrus.WithField("ref", "api").
				Errorf("create index %q: not acknowledged", store.ElasticSearchIndexName)
		}

		// OK
	}

	store.ESClient = db

	return nil
}
