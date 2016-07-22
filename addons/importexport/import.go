package importexport

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/store"
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

			if err := bucket.Import(_bucket.Data); err != nil {
				logrus.Errorf("Import data bucket %q, %#v", bucket.Name(), err)
			}

			bucket.InitRawDataStore(bucket.GetRawDataStoreType(),
				bucket.GetRawDataStoreNameWithoutPostfix())
			bucket.InitMetaDataStore(bucket.GetMetaDataStoreType(),
				bucket.GetMetaDataStoreNameWithoutPostfix())
			bucket.InitMapDataStore(bucket.GetMapDataStoreType(),
				bucket.GetMapDataStoreNameWithoutPostfix())

			if err := bucket.UpdateMapping(); err != nil {
				logrus.WithField("ref", addonName).
					Errorf("Update mapping bucket %q, %#v", bucket.Name(), err)
			}

			if err := bucket.Sync(); err != nil {
				logrus.WithField("ref", addonName).
					Errorf("Save bucket %q, %#v", bucket.Name(), err)
			}
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

		if err := file.Import(_file.Data); err != nil {
			logrus.WithField("ref", addonName).
				Errorf("Import data file %q:%q, %v", _file.Name, _file.Bucket, err)
		}
		if err := file.Sync(); err != nil {
			logrus.WithField("ref", addonName).
				Errorf("Save file %q:%q, %v", _file.Name, _file.Bucket, err)
		}
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
