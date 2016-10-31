package api

// InitFirstRun installation of first run
func InitFirstRunIfNeed() error {
	if !IsFirstStart() {
		return nil
	}

	logger.Println("init of first run ...")

	/*
	   загружаем архив из интернет
	   устанавливаем его
	*/

	return nil
}
