package importexport

import (
	"api/config"
	apiutils "api/utils"
	"utils"
)

func MainSettings() *Settings {
	return config.Cfgx.Config(addonName).(*Settings)
}

type Settings struct {
	// same value as `addonName`
	*settings `toml:"importexport" json:"importexport"`
}

func (s Settings) Merge(cfg interface{}) error {
	return apiutils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	Groups map[string][]GroupSettings `toml:"groups" json:"groups"`
}

func (s Settings) GroupNames() (arr []string) {
	for name, _ := range s.Groups {
		arr = append(arr, name)
	}

	return arr
}

func (s Settings) SettingsForBucket(groupName, bucketName string) *GroupSettings {
	for _, setting := range s.Groups[groupName] {
		if setting.BucketName == bucketName {
			return &setting
		}
	}

	return nil
}

type GroupSettings struct {
	BucketName string   `toml:"bucket" json:"bucket"`
	Files      []string `toml:"files" json:"files"`
	All        bool     `toml:"all" json:"all"`
}

func (gs GroupSettings) IncludeFile(filename string) bool {
	return utils.NewA(gs.Files).Include(filename)
}
