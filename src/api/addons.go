package api

import (
	"addons"
	"addons/example"
)

func SetupAddons() error {
	addons.Addons[example.NAME] = example.NewAddon()

	return nil
}
