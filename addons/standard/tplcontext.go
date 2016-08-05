package standard

import (
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/inpime/dbox"
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/store"
	"github.com/inpime/sdata"
)

func (Extension) initTplContext() {
	pongo2.DefaultSet.Globals["NewUUID"] = func() *pongo2.Value {
		return pongo2.AsValue(dbox.NewUUID())
	}

	pongo2.DefaultSet.Globals["Config"] = func() *pongo2.Value {
		return pongo2.AsValue(config.Cfgx.Config(addonName).(*Settings).Config)
	}

	pongo2.DefaultSet.Globals["SectionAppConfig"] = func(sectionName *pongo2.Value) *pongo2.Value {
		return pongo2.AsValue(config.Cfgx.Config(sectionName.String()))
	}

	// DeleteFile
	pongo2.DefaultSet.Globals["DeleteFile"] = func(bucketId, fileId *pongo2.Value) *pongo2.Value {
		if !bucketId.IsString() || !fileId.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		return pongo2.AsValue(store.DeleteFileID(bucketId.String(), fileId.String()))
	}

	pongo2.DefaultSet.Globals["NewFile"] = func(bucketName *pongo2.Value) *pongo2.Value {
		if !bucketName.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, _ := store.LoadOrNewFileID(bucketName.String(), "")

		return pongo2.AsValue(file)
	}

	// LoadByID load file by ID
	pongo2.DefaultSet.Globals["LoadByID"] = func(
		bucketName,
		fileID *pongo2.Value,
	) *pongo2.Value {

		if !bucketName.IsString() || !fileID.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, err := store.LoadOrNewFileID(
			strings.ToLower(bucketName.String()),
			fileID.String())

		if err != nil {
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(file)
	}

	// Load load file by name
	pongo2.DefaultSet.Globals["Load"] = func(
		bucketName,
		fileName *pongo2.Value,
	) *pongo2.Value {

		if !bucketName.IsString() || !fileName.IsString() {
			return pongo2.AsValue(ErrNotValidData)
		}

		file, _ := store.LoadOrNewFile(
			strings.ToLower(bucketName.String()),
			fileName.String())

		return pongo2.AsValue(file)
	}

	//

	pongo2.DefaultSet.Globals["URLQuery"] = func(args ...*pongo2.Value) *pongo2.Value {
		emptyUrl, _ := url.Parse("")

		if len(args) == 0 {
			return pongo2.AsValue(emptyUrl)
		}

		_url, ok := args[0].Interface().(*url.URL)

		if !ok {
			logrus.WithFields(logrus.Fields{
				"_api": addonName,
			}).Warningf("not expected type %T, want '*url.URL'", args[0].Interface())

			return pongo2.AsValue(emptyUrl)
		}

		if (len(args)-1)%2 != 0 {
			logrus.WithFields(logrus.Fields{
				"_api": addonName,
			}).Warningf("args expected in multiples of two, want %d", len(args)-1)
			return pongo2.AsValue(emptyUrl)
		}

		queryParams := args[1:]
		urlQueryValues := _url.Query()

		for i := 0; i < len(queryParams); i += 2 {
			urlQueryValues.Add(queryParams[i].String(), queryParams[i+1].String())
		}

		_url.RawQuery = urlQueryValues.Encode()

		return pongo2.AsValue(_url)
	}

	pongo2.DefaultSet.Globals["M"] = func() *pongo2.Value {

		return pongo2.AsValue(sdata.NewStringMap())
	}

	pongo2.DefaultSet.Globals["A"] = func() *pongo2.Value {
		return pongo2.AsValue(sdata.NewArray())
	}

	pongo2.DefaultSet.Globals["Validator"] = func() *pongo2.Value {

		return pongo2.AsValue(NewValidatorData())
	}

	// alias Validator
	pongo2.DefaultSet.Globals["V"] = func() *pongo2.Value {

		return pongo2.AsValue(NewValidatorData())
	}

	// CreateBucket special function (used only to create a bucket)
	pongo2.DefaultSet.Globals["CreateBucket"] = func(_opt *pongo2.Value) *pongo2.Value {
		opt := sdata.NewStringMapFrom(_opt.Interface())

		bucketName := opt.String("Name")
		bucket, err := store.BucketByName(bucketName)

		if opt.Bool("SameAsMetaStoreType") {
			bucket.InitInOneStore(dbox.StoreType(opt.String("MetaDataStoreType")))
		} else {

			bucket.InitMetaDataStore(
				dbox.StoreType(opt.String("MetaDataStoreType")),
				opt.Bool("MetaHaveSuffix")) // store key - <type>.<name>.meta
			bucket.InitMapDataStore(
				dbox.StoreType(opt.String("MapDataStoreType")),
				opt.Bool("MapDataHaveSuffix")) // store key - <type>.<name>.mapdata
			bucket.InitRawDataStore(
				dbox.StoreType(opt.String("RawDataStoreType")),
				opt.Bool("RawDataHaveSuffix")) // store key - <type>.<name>.rawdata
		}

		// Only create new bucket
		if err != dbox.ErrNotFound {
			return pongo2.AsValue(err)
		}

		bucket.MMetaDataFilesMapping().LoadFrom(store.FileMetaMappingDefault)
		bucket.MMapDataFilesMapping().LoadFrom(opt.String("MappingMapDataFiles"))

		if err := bucket.UpdateMapping(); err != nil {
			logrus.WithError(err).Errorf("create new bucket %q: update mapping", bucketName)
			return pongo2.AsValue(err)
		}

		if err := bucket.Sync(); err != nil {
			logrus.WithError(err).Errorf("create new bucket %q: save bucket", bucketName)
			return pongo2.AsValue(err)
		}

		return pongo2.AsValue(bucket)
	}

	pongo2.DefaultSet.Globals["SendEmail"] = func(to, subject, template, context *pongo2.Value) *pongo2.Value {
		trackid, err := SendEmail(to.String(), subject.String(), template.String(), context.Interface())

		if err != nil {
			logrus.WithError(err).
				WithFields(logrus.Fields{
					"_api":     addonName,
					"to":       to.String(),
					"subject":  subject.String(),
					"template": template.String(),
					"context":  context.Interface(),
				}).Error("Send email")
			return pongo2.AsValue(err)
		}

		logrus.
			WithFields(logrus.Fields{
				"_api":     addonName,
				"to":       to.String(),
				"subject":  subject.String(),
				"template": template.String(),
				"context":  context.Interface(),
				"trackid":  trackid,
			}).Info("Send email")

		return pongo2.AsValue(trackid)
	}
}
