package importexport

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
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
	config := MainSettings().SettingsForBucket(groupName, bucketName)

	if config == nil || config.BucketName == "" {
		return false
	}

	if config.All {
		return true
	}

	return len(config.Files) > 0
}

// IsIncludeInGroupFileImportExport является ли файл системным согласно настройкам
func IsIncludeInGroupFileImportExport(groupName, bucketName, fileName string) bool {
	if !IsIncludeInGroupBucketImportExport(groupName, bucketName) {
		return false
	}

	config := MainSettings().SettingsForBucket(groupName, bucketName)

	// не проверяем потому что выше IsIncludeInGroupBucketImportExport

	if config.All {
		return true
	}

	return config.IncludeFile(fileName)
}

// ListGroupsImportExport возвращает список групп указанных в настройках приложения
func ListGroupsImportExport() []string {
	return MainSettings().GroupNames()
}
