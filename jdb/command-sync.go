package jdb

func (s *Command) sync() error {
	if err := s.prepare(); err != nil {
		return err
	}
	model := s.From

	if len(s.Data) == 0 {
		return nil
	}

	id := s.Data[0].Str(SYSID)
	if id == "" {
		return nil
	}

	item, err := model.
		Update(s.Data[0]).
		Where(SYSID).Eq(id).
		setSync().
		OneTx(s.tx)
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	_, err = model.
		Insert(s.Data[0]).
		setSync().
		ExecTx(s.tx)
	if err != nil {
		return err
	}

	return nil
}
