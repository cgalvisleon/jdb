package jdb

import (
	"database/sql"

	"github.com/cgalvisleon/et/utility"
)

type Tx struct {
	Id        string
	Committed bool
	Tx        *sql.Tx
}

/**
* NewTx
* @return *Tx
**/
func NewTx() *Tx {
	return &Tx{
		Id: utility.UUID(),
	}
}

/**
* InitTx
* @param tx *Tx
* @return *Tx
**/
func InitTx(tx *Tx) *Tx {
	if tx.Tx == nil {
		tx = NewTx()
	}

	return tx
}

/**
* Begin
* @param db *sql.DB
* @return error
**/
func (s *Tx) Begin(db *sql.DB) error {
	if s.Tx != nil {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	s.Tx = tx

	return nil
}

/**
* Commit
* @return error
**/
func (s *Tx) Commit() error {
	if s.Tx == nil {
		return nil
	}

	if s.Committed {
		return nil
	}

	err := s.Tx.Commit()
	s.Committed = true

	return err
}

/**
* Rollback
* @return error
**/
func (s *Tx) Rollback() error {
	if s.Tx == nil {
		return nil
	}

	if s.Committed {
		return nil
	}

	err := s.Tx.Rollback()
	s.Committed = true

	return err
}
