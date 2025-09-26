package jdb

/**
* initCore
* @return error
**/
func initCore(db *Database) error {
	if err := defineSeries(db); err != nil {
		return err
	}
	if err := defineModel(db); err != nil {
		return err
	}

	return nil
}
