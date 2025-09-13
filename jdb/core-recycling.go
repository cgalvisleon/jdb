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
	if s.driver.Name() == SqliteDriver {
		return nil
	}

	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreRecycling != nil {
		return nil
	}

	coreRecycling = NewModel(coreSchema, "recycling", 1)
	coreRecycling.DefineColumn(cf.CreatedAt, TypeDataDateTime)
	coreRecycling.DefineColumn(cf.UpdatedAt, TypeDataDateTime)
	coreRecycling.DefineColumn("schema_name", TypeDataText)
	coreRecycling.DefineColumn("table_name", TypeDataText)
	coreRecycling.DefineColumn(cf.SystemId, TypeDataKey)
	coreRecycling.DefineIndexField()
	coreRecycling.DefinePrimaryKey("schema_name", "table_name", cf.SystemId)
	coreRecycling.DefineIndex(true,
		cf.Index,
	)
	if err := coreRecycling.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* upsertRecycling
* @param tx *Tx, schema, name, sysId string
* @return error
**/
func (s *DB) upsertRecycling(tx *Tx, schema, name, sysId string) error {
	if coreRecycling == nil || !coreRecycling.isInit {
		return nil
	}

	now := timezone.Now()
	data := et.Json{
		"schema_name": schema,
		"table_name":  name,
		cf.SystemId:   sysId,
	}
	_, err := coreRecycling.
		Upsert(data).
		BeforeInsert(func(tx *Tx, data et.Json) error {
			data.Set(cf.CreatedAt, now)
			data.Set(cf.UpdatedAt, now)
			return nil
		}).
		BeforeUpdate(func(tx *Tx, data et.Json) error {
			data.Set(cf.UpdatedAt, now)
			return nil
		}).
		ExecTx(tx)
	if err != nil {
		return err
	}

	return nil
}

/**
* GetRecycling
* @param schema, name, sysId string
* @return et.Item, error
**/
func (s *DB) GetRecycling(schema, name, sysId string) (et.Item, error) {
	if coreRecycling == nil || !coreRecycling.isInit {
		return et.Item{}, mistake.New(MSG_DATABASE_NOT_CONCURRENT)
	}

	if !utility.ValidName(schema) {
		return et.Item{}, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "schema_name")
	}

	if !utility.ValidName(name) {
		return et.Item{}, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "table_name")
	}

	if !utility.ValidId(sysId) {
		return et.Item{}, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, cf.SystemId)
	}

	item, err := coreRecycling.
		Where("schema_name").Eq(schema).
		And("table_name").Eq(name).
		And(cf.SystemId).Eq(sysId).
		One()
	if err != nil {
		return et.Item{}, err
	}

	return item, nil
}

/**
* deleteRecycling
* @param tx *Tx, schema, name, sysId string
* @return error
**/
func (s *DB) deleteRecycling(tx *Tx, schema, name, sysId string) error {
	if coreRecycling == nil || !coreRecycling.isInit {
		return nil
	}

	if !utility.ValidName(schema) {
		return mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "schema_name")
	}

	if !utility.ValidName(name) {
		return mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "table_name")
	}

	item, err := coreRecycling.
		Delete().
		Where("schema_name").Eq(schema).
		And("table_name").Eq(name).
		And(cf.SystemId).Eq(sysId).
		ExecTx(tx)
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
	if coreRecycling == nil || !coreRecycling.isInit {
		return et.Item{}, mistake.New(MSG_DATABASE_NOT_CONCURRENT)
	}

	result, err := coreRecycling.
		Query(search)
	if err != nil {
		return et.List{}, err
	}

	return result, nil
}

/**
* HttpGetRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HttpGetRecycling(w http.ResponseWriter, r *http.Request) {
	schema := r.PathValue("schema")
	name := r.PathValue("name")
	id := r.PathValue("id")
	result, err := s.GetRecycling(schema, name, id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, result)
}

/**
* HttpDeleteRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HttpDeleteRecycling(w http.ResponseWriter, r *http.Request) {
	schema := r.PathValue("schema")
	name := r.PathValue("name")
	id := r.PathValue("id")
	err := s.deleteRecycling(nil, schema, name, id)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"message": "Recycling deleted successfully",
			"schema":  schema,
			"name":    name,
			"id":      id,
		},
	})
}

/**
* HttpQueryRecycling
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HttpQueryRecycling(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := s.QueryRecycling(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RESULT(w, r, http.StatusOK, result)
}
