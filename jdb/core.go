package jdb

/**
* initCore
* @return error
**/
func initCore(db *DB) error {
	if err := defineModel(db); err != nil {
		return err
	}
	if err := defineRecords(db); err != nil {
		return err
	}
	if err := defineSeries(db); err != nil {
		return err
	}

	return nil
}
