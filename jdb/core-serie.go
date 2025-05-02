package jdb

import (
	"net/http"

	"github.com/cgalvisleon/et/console"
	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/response"
	"github.com/cgalvisleon/et/strs"
	"github.com/cgalvisleon/et/timezone"
	"github.com/cgalvisleon/et/utility"
)

var coreSeries *Model

func (s *DB) defineSeries() error {
	if err := s.defineSchema(); err != nil {
		return err
	}

	if coreSeries != nil {
		return nil
	}

	coreSeries = NewModel(coreSchema, "series", 1)
	coreSeries.DefineColumn(CREATED_AT, CreatedAtField.TypeData())
	coreSeries.DefineColumn(UPDATED_AT, UpdatedAtField.TypeData())
	coreSeries.DefineColumn("tag", TypeDataText)
	coreSeries.DefineColumn("value", TypeDataInt)
	coreSeries.DefineIndexField()
	coreSeries.DefinePrimaryKey("tag")
	coreSeries.DefineIndex(true,
		"tag",
		INDEX,
	)
	if err := coreSeries.Init(); err != nil {
		return console.Panic(err)
	}

	return nil
}

/**
* CurrentSerie
* @param tag string
* @return int64, error
**/
func (s *DB) CurrentSerie(tag string) (int64, error) {
	if !utility.ValidStr(tag, 0, []string{}) {
		return 0, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "tag")
	}

	if !s.UseCore {
		return 0, nil
	}

	item, err := coreSeries.
		Where("tag").Eq(tag).
		One()
	if err != nil {
		return 0, err
	}

	if !item.Ok {
		return 0, mistake.Newf(MSG_SERIE_NOT_FOUND, tag)
	}

	return item.Int64("value"), nil
}

/**
* GetSerie
* @param tag string
* @return int64, error
**/
func (s *DB) GetSerie(tag string) (int64, error) {
	if !utility.ValidStr(tag, 0, []string{}) {
		return 0, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "tag")
	}

	now := timezone.Now()
	_, err := coreSeries.
		Upsert(et.Json{
			CREATED_AT: now,
			UPDATED_AT: now,
			"tag":      tag,
			"value":    1,
		}).
		Return("value").
		Exec()
	if err != nil {
		return 0, err
	}

	return 1, nil
}

/**
* GetCode
* @param tag, format string
* @return string, error
**/
func (s *DB) GetCode(tag, format string) (string, error) {
	num, err := s.GetSerie(tag)
	if err != nil {
		return "", err
	}

	if len(format) == 0 {
		return strs.Format("%08v", num), nil
	} else {
		return strs.FormatSerie(format, num), nil
	}
}

/**
* SetSerie
* @param tag string, val int64
* @return int64, error
**/
func (s *DB) SetSerie(tag string, val int64) (int64, error) {
	if !utility.ValidStr(tag, 0, []string{}) {
		return 0, mistake.Newf(MSG_ATTRIBUTE_REQUIRED, "tag")
	}

	now := timezone.Now()
	item, err := coreSeries.
		Update(et.Json{
			UPDATED_AT: now,
			"value":    val,
		}).
		Where("tag").Eq(tag).
		Return("value").
		One()
	if err != nil {
		return 0, err
	}

	if item.Ok {
		return item.Int64("value"), nil
	}

	_, err = coreSeries.
		Insert(et.Json{
			CREATED_AT: now,
			UPDATED_AT: now,
			"tag":      tag,
			"value":    val,
		}).
		Return("value").
		Exec()
	if err != nil {
		return 0, err
	}

	return val, nil
}

/**
* DeleteSerie
* @param tag string
* @return error
**/
func (s *DB) DeleteSerie(tag string) error {
	item, err := coreSeries.
		Delete("tag").Eq(tag).
		Exec()
	if err != nil {
		return err
	}

	if !item.Ok {
		return mistake.Newf(MSG_SERIE_NOT_FOUND, tag)
	}

	return nil
}

/**
* QuerySerie
* @param search et.Json
* @return interface{}, error
**/
func (s *DB) QuerySerie(search et.Json) (interface{}, error) {
	result, err := coreSeries.
		Query(search)
	if err != nil {
		return et.List{}, err
	}

	return result, nil
}

/**
* HandlerCurrentSerie
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerCurrentSerie(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	serie, err := s.CurrentSerie(tag)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"tag":   tag,
			"value": serie,
		},
	})
}

/**
* HandlerGetSeries
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerGetSeries(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	serie, err := s.GetSerie(tag)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"tag":   tag,
			"value": serie,
		},
	})
}

/**
* HandlerGetCode
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerGetCode(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	format := r.URL.Query().Get("format")
	code, err := s.GetCode(tag, format)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"tag":    tag,
			"format": format,
			"code":   code,
		},
	})
}

/**
* HandlerSetSerie
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerSetSeries(w http.ResponseWriter, r *http.Request) {
	params, _ := response.GetBody(r)
	tag := params.String("tag")
	val := params.Int64("value")
	_, err := s.SetSerie(tag, val)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"tag":   tag,
			"value": val,
		},
	})
}

/**
* HandlerQuerySerie
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerQuerySeries(w http.ResponseWriter, r *http.Request) {
	body, _ := response.GetBody(r)
	result, err := s.QuerySerie(body)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.RESULT(w, r, http.StatusOK, result)
}

/**
* HandlerDeleteSeries
* @param w http.ResponseWriter
* @param r *http.Request
**/
func (s *DB) HandlerDeleteSeries(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	err := s.DeleteSerie(tag)
	if err != nil {
		response.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.ITEM(w, r, http.StatusOK, et.Item{
		Ok: true,
		Result: et.Json{
			"message": "Serie deleted successfully",
			"tag":     tag,
		},
	})
}
