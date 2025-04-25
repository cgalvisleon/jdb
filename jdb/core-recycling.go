package jdb

import (
	"net/http"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

var coreRecycling *Model

func (s *DB) defineRecycling() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreRecycling != nil {
		return nil
	}

	coreRecycling = NewModel(coreSchema, "recycling", 1)
	coreRecycling.DefineColumn(CREATED_AT, CreatedAtField.TypeData())
	coreRecycling.DefineColumn(UPDATED_AT, UpdatedAtField.TypeData())
	coreRecycling.DefineColumn("table_name", TypeDataText)
	coreRecycling.DefineColumn(SYSID, SystemKeyField.TypeData())
	coreRecycling.DefineIndexField()
	coreRecycling.DefinePrimaryKey("table_name", SYSID)
	coreRecycling.DefineIndex(true,
		INDEX,
	)
	if err := coreRecycling.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* upsertRecycling
* @param table, sysId, statusId string
* @return error
**/
func (s *DB) upsertRecycling(table, sysId, statusId string) error {
	if statusId != utility.FOR_DELETE {
		_, err := coreRecycling.
			Delete().
			Where("table_name").Eq(table).
			And(SYSID).Eq(sysId).
			Exec()
		if err != nil {
			return err
		}

		return nil
	}

	now := timezone.Now()
	item, err := coreRecycling.
		Update(et.Json{
			UPDATED_AT: now,
		}).
		Where("table_name").Eq(table).
		And(SYSID).Eq(sysId).
		One()
	if err != nil {
		return err
	}

	if item.Ok {
		return nil
	}

	_, err = coreRecycling.
		Insert(et.Json{
			CREATED_AT:   now,
			UPDATED_AT:   now,
			"table_name": table,
			SYSID:        sysId,
		}).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

/**
* GetRecycling
* @param table, sysId string
* @return et.Item, error
**/
func (s *DB) GetRecycling(table, sysId string) (et.Item, error) {
	if !utility.ValidName(table) {
		return et.Item{}, mistake.Newf(MSG_ATTR_REQUIRED, "table")
	}

	if !utility.ValidId(sysId) {
		return et.Item{}, mistake.Newf(MSG_ATTR_REQUIRED, SYSID)
	}

	item, err := coreRecycling.
		Where("table_name").Eq(table).
		And(SYSID).Eq(sysId).
		One()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* deleteRecycling
* @param table, sysId string
* @return error
**/
func (s *DB) deleteRecycling(table, sysId string) error {
	if !utility.ValidName(table) {
		return mistake.Newf(MSG_ATTR_REQUIRED, "table")
	}

	item, err := coreRecycling.
		Delete().
		Where("table_name").Eq(table).
		And(SYSID).Eq(sysId).
		One()
	if err != nil {
		return err
	}

	if !item.Ok {
		return mistake.New(MSG_RECORD_NOT_FOUND)
	}

	return nil
}

/**
* QueryRecycling
* @param search et.Json
* @return interface{}, error
**/
func (s *DB) QueryRecycling(search et.Json) (interface{}, error) {
	result, err := coreRecycling.
		Query(search)
	if err != nil {
		return et.List{}, err
	}

	return result, nil
}

/**
* HandlerGetRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerGetRecycling(w http.ResponseWriter, r *http.Request) {
	table := r.PathValue("table")
	id := r.PathValue("id")
	result, err := s.GetRecycling(table, id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* HandlerDeleteRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerDeleteRecycling(w http.ResponseWriter, r *http.Request) {
	table := r.PathValue("table")
	id := r.PathValue("id")
	err := s.deleteRecycling(table, id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"message": "Recycling deleted successfully",
			"table":   table,
			"id":      id,
		},
	})
}

/**
* HandlerQueryRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerQueryRecycling(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := s.QueryRecycling(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RESULT(w, r, http.StatusOK, result)
}
