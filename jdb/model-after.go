package jdb

import (
	"github.com/cgalvisleon/et/et"
)

/**
* afterInsertDefault
* @param tx *Tx
* @param data et.Json
* @return error
**/
func (s *Model) afterInsertDefault(tx *Tx, data et.Json) error {
	if s.UseCore {
		s.Db.upsertTable(tx, s.Schema, s.Name, 1)

		if s.SystemKeyField != nil {
			sysId := data.Str(s.SystemKeyField.Name)
			err := s.Db.upsertRecord(tx, s.Schema, s.Name, sysId, "insert")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/**
* afterUpdateDefault
* @param tx *Tx
* @param data et.Json
* @return error
**/
func (s *Model) afterUpdateDefault(tx *Tx, data et.Json) error {
	if s.UseCore && s.SystemKeyField != nil {
		sysId := data.Str(s.SystemKeyField.Name)
		err := s.Db.upsertRecord(tx, s.Schema, s.Name, sysId, "update")
		if err != nil {
			return err
		}

		if s.StatusField != nil {
			oldStatus := data.ValStr(ACTIVE, s.StatusField.Name)
			newStatus := data.ValStr(oldStatus, s.StatusField.Name)
			if oldStatus != newStatus && newStatus == FOR_DELETE {
				err := s.Db.upsertRecycling(tx, s.Schema, s.Name, sysId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/**
* afterDeleteDefault
* @param tx *Tx
* @param data et.Json
* @return error
**/
func (s *Model) afterDeleteDefault(tx *Tx, data et.Json) error {
	if s.UseCore {
		s.Db.upsertTable(tx, s.Schema, s.Name, -1)

		if s.SystemKeyField != nil {
			sysId := data.Str(s.SystemKeyField.Name)
			err := s.Db.upsertRecord(tx, s.Schema, s.Name, sysId, "delete")
			if err != nil {
				return err
			}

			if s.StatusField != nil {
				err := s.Db.deleteRecycling(tx, s.Schema, s.Name, sysId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/**
* AfterInsert
* @param fn DataFunction
* @return *Command
**/
func (s *Model) AfterInsert(fn DataFunctionTx) *Model {
	s.afterInsert = append(s.afterInsert, fn)

	return s
}

/**
* AfterUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterUpdate(fn DataFunctionTx) *Model {
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}

/**
* AfterDelete
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterDelete(fn DataFunctionTx) *Model {
	s.afterDelete = append(s.afterDelete, fn)

	return s
}

/**
* AfterInsertOrUpdate
* @param fn DataFunctionTx
* @return *Command
**/
func (s *Model) AfterInsertOrUpdate(fn DataFunctionTx) *Model {
	s.afterInsert = append(s.afterInsert, fn)
	s.afterUpdate = append(s.afterUpdate, fn)

	return s
}
