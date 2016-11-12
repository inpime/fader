package api

import (
	"strings"
	"errors"
	"io/ioutil"
	"interfaces"
)

// InitFirstRunIfNeed installation of first run
func InitFirstRunIfNeed() error {
	logger.Println("init of first run ...")

	if !IsFirstStart() {
		logger.Println("this is not the first run... skiped")
		return nil
	}

	var data []byte 

	if strings.HasPrefix("http", settings.InitFile) {
		logger.Println("download from the internet is not implemented")
		return errors.New("not implemented")
	}

	data, err := ioutil.ReadFile(settings.InitFile)

	if err != nil {
		logger.Println("error open file", err)
		return errors.New("internat error")
	}

	importer := interfaces.NewImportManager(
		bucketManager,
		fileManager,
	)

	info, err := importer.Import(data)

	if err != nil {
		logger.Println("error import file", err)
		return errors.New("internal error")
	}

	logger.Println("the imported file, appname -", info.AppName())
	logger.Println("the imported file, version -", info.Version())
	logger.Println("the imported file, author -", info.Author())
	logger.Println("the imported file, description -", info.Description())
	logger.Println("the imported file, date -", info.DateTime())

	logger.Println("import completed")

	return nil
}