package jdb

func (s *Command) relationsTo() error {
	for _, relationsTo := range s.RelationsTo {
		for _, field := range relationsTo {
			detail := field.Column.Detail
			if detail == nil {
				continue
			}

			with := detail.With
			if with == nil {
				continue
			}

			val, err := field.valueJson()
			if err != nil {
				return err
			}

			_, err = with.
				Upsert(val).
				setDebug(s.IsDebug).
				ExecTx(s.tx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
