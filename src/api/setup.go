package api

import (
	"errors"
	"interfaces"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	configBucketName   = "settings"
	mainConfigFileName = "main.toml"

	appConfigUpdateFn = func() error {
		if config == nil {
			logger.Println("[RefreshAppConfig] app config is empty")
			return errors.New("invalid data")
		}

		_config := newConfig()

		mainTomlFile, err := fileManager.FindFileByName(
			configBucketName, mainConfigFileName,
			interfaces.FullFile,
		)

		if err != nil {
			logger.Println("[RefreshAppConfig] find file ", mainConfigFileName, ":", err)
			// TODO: signal for setup application
			return err
		}

		if _, err := toml.Decode(string(mainTomlFile.RawData), _config); err != nil {
			logger.Println("[RefreshAppConfig] decode toml, ", err)
			return err
		}

		for _, includeFileName := range _config.Config().Main.Include {
			__config := newConfig()

			includeTomlFile, err := fileManager.FindFileByName(
				configBucketName, includeFileName,
				interfaces.FullFile,
			)

			if err != nil {
				logger.Println("[RefreshAppConfig] find include file ", includeFileName, ":", err)
				// TODO: signal for setup application
				return err
			}

			if _, err := toml.Decode(string(includeTomlFile.RawData), __config); err != nil {
				logger.Println("[RefreshAppConfig] decode toml, ", err)
				return err
			}

			if err := _config.Merge(__config); err != nil {
				logger.Println("[RefreshAppConfig] merge config file ", includeFileName, ":", err)
				return err
			}
		}

		config.Lock()
		config.Main = _config.Main
		config.Routing = _config.Routing
		config.Addons = _config.Addons
		config.Unlock()

		return nil
	}

	appRoutesUpdateFn = func() error {

		vrouter.RefreshRoutes(config.Config().Routing.Routs)

		return nil
	}
)

func InitAppConfigurator() {

	logger.Println("[AppConfigurator] start periodic updates app config...")
	RefreshEvery(time.Second*1, appConfigUpdateFn)
}
