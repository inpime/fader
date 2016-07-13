package importexport

import (
	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/store"
	"time"
)

// makeImportImportExport выполнить импорт из архива
func makeImportImportExport(data []byte) error {
	archive := newArchivePkg()
	err := archive.Import(data)

	if err != nil {
		return err
	}

	for _, _bucket := range archive.Buckets {
		bucket, err := store.BucketByName(_bucket.Name)

		if err == dbox.ErrNotFound {
			logrus.WithField("_api", addonName).
				Infof("import: Create a bucket %q", _bucket.Name)

			bucket.SetID(_bucket.ID)
			bucket.SetName(_bucket.Name)
			bucket.SetBucket(_bucket.Bucket)
			bucket.Import(_bucket.Data)

			bucket.InitRawDataStore(bucket.GetRawDataStoreType(),
				bucket.GetRawDataStoreNameWithoutPostfix())
			bucket.InitMetaDataStore(bucket.GetMetaDataStoreType(),
				bucket.GetMetaDataStoreNameWithoutPostfix())
			bucket.InitMapDataStore(bucket.GetMapDataStoreType(),
				bucket.GetMapDataStoreNameWithoutPostfix())

			bucket.UpdateMapping()
			bucket.Sync()
		}
	}

	time.Sleep(time.Second * 3)

	for _, _file := range archive.Files {
		file, err := store.LoadOrNewFile(_file.Bucket, _file.Name)

		if err == dbox.ErrNotFound {
			file.SetID(_file.ID)
			file.SetName(_file.Name)
			file.SetBucket(_file.Bucket)
		}

		logrus.WithField("_api", addonName).
			Infof("import: Upsert a file %q", _file.Bucket+"@"+_file.Name)

		file.Import(_file.Data)
		file.Sync()
	}

	return nil
}

// AppImportFromLastArchive загрузить последнюю версию консоль панели и загрузить ее
func AppImportFromLastArchive(archive string) error {

	logrus.WithField("_api", addonName).Infof("Downloading on link %q...", archive)

	data, err := loadLatestArchive(archive)

	if err != nil {
		logrus.WithField("_api", addonName).Errorf("Error download, more details: %q", err)
		return err
	}

	if err := makeImportImportExport(data); err != nil {
		logrus.WithField("_api", addonName).WithError(err).Errorf("Error importing.")
		return err
	}

	logrus.WithField("_api", addonName).Debug("OK importing.")

	return nil
}
