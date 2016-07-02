package importexport

import (
	"api/config"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"utils"
)

// loadLatestArchive load latest version archive
func loadLatestArchive(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithField("_service", addonName).
			Errorf("not successful download %q, %q", url, resp.Status)
		return []byte{}, fmt.Errorf("not successful request, got %q, want `200 OK`", resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// IsIncludeInGroupBucketImportExport является ли файл системным согласно настройкам
func IsIncludeInGroupBucketImportExport(groupName, bucketName string) bool {
	config := config.AppSettings().M(SettingsSectionNameKey).M(groupName)

	if !config.Include(bucketName) {
		return false
	}

	bucketConfig := config.M(bucketName)

	if bucketConfig.Bool("all") {
		return true
	}

	bucketFiles := utils.NewA(bucketConfig.Strings("files"))

	return bucketFiles.Len() > 0
}

// IsIncludeInGroupFileImportExport является ли файл системным согласно настройкам
func IsIncludeInGroupFileImportExport(groupName, bucketName, fileName string) bool {
	if !IsIncludeInGroupBucketImportExport(groupName, bucketName) {
		return false
	}

	config := config.AppSettings().M(SettingsSectionNameKey).M(groupName)

	bucketConfig := config.M(bucketName)

	if bucketConfig.Bool("all") {
		return true
	}

	bucketFiles := utils.NewA(bucketConfig.Strings("files"))

	return bucketFiles.Include(fileName)
}

// ListGroupsImportExport возвращает список групп указанных в настройках приложения
func ListGroupsImportExport() []string {
	return config.AppSettings().Keys(SettingsSectionNameKey)
}
