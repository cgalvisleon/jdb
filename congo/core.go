package jdb

/**
* initCore
* @return error
**/
func initCore() error {
	if err := defineModel(); err != nil {
		return err
	}
	if err := defineRecord(); err != nil {
		return err
	}
	if err := defineRecycling(); err != nil {
		return err
	}
	if err := defineSeries(); err != nil {
		return err
	}
	if err := defineTables(); err != nil {
		return err
	}

	return nil
}
