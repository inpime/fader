package importexport

import (
	// "github.com/Sirupsen/logrus"
	"github.com/inpime/fader/api/config"
	"github.com/inpime/fader/utils"
	"github.com/inpime/fader/utils/sdata"
)

func MainSettings() *Settings {
	return config.Cfgx.Config(addonName).(*Settings)
}

type Settings struct {
	// same value as `addonName`
	*settings `toml:"importexport" json:"importexport"`
}

func (s Settings) Merge(cfg interface{}) error {
	return utils.AppendOrReplace(s.settings, cfg.(*Settings).settings)
}

type settings struct {
	Groups []GroupSettings `toml:"groups" json:"groups"`
}

func (s Settings) GroupNames() (arr []string) {
	groups := sdata.NewStringMap()

	for _, group := range s.Groups {
		if _, exist := groups.GetIf(group.GroupName); !exist {
			groups.Set(group.GroupName, nil)
		}
	}

	return groups.Keys()
}

func (s Settings) SettingsForBucket(groupName, bucketName string) *GroupSettings {
	for _, group := range s.Groups {
		if group.GroupName == groupName &&
			group.BucketName == bucketName {
			return &group
		}
	}

	return nil
}

type GroupSettings struct {
	GroupName  string   `toml:"group" json:"group"`
	BucketName string   `toml:"bucket" json:"bucket"`
	Files      []string `toml:"files" json:"files"`
	All        bool     `toml:"all" json:"all"`
}

func (gs GroupSettings) IncludeFile(filename string) bool {
	return sdata.NewArrayFrom(gs.Files).Index(filename) > -1
}
