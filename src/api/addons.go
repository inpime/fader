package api

import (
	"addons"
	"addons/example"
)

func SetupAddons() error {
	addons.Addons[example.NAME] = example.NewAddon()
	addons.Addons[ADDONS_BASIC_NAME] = NewBasicAddon()

	return nil
}
