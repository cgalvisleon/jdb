package jdb

import (
	"slices"
	"strconv"
	"strings"

	"github.com/cgalvisleon/et/et"
	"github.com/cgalvisleon/et/mistake"
	"github.com/cgalvisleon/et/strs"
)

type TypeDefinition int

const (
	TypeDefinitionColumn TypeDefinition = iota
	TypeDefinitionIndex
	TypeDefinitionUnique
	TypeDefinitionPrimaryKey
	TypeDefinitionPrimaryKeyField
	TypeDefinitionForeignKey
	TypeDefinitionSourceField
	TypeDefinitionAtribute
	TypeDefinitionCreatedAtField
	TypeDefinitionUpdatedAtField
	TypeDefinitionStatusField
	TypeDefinitionSystemKeyField
	TypeDefinitionIndexField
	TypeDefinitionFullTextField
	TypeDefinitionProjectField
	TypeDefinitionHidden
	TypeDefinitionRequired
	TypeDefinitionRelation
	TypeDefinitionRollup
	TypeDefinitionDetail
	TypeDefinitionHistory
	TypeDefinitionModel
	TypeDefinitionProjectModel
)

func (s TypeDefinition) Str() string {
	switch s {
	case TypeDefinitionColumn:
		return "column"
	case TypeDefinitionIndex:
		return "index"
	case TypeDefinitionUnique:
		return "unique"
	case TypeDefinitionPrimaryKey:
		return "primary_key"
	case TypeDefinitionPrimaryKeyField:
		return "primary_key_field"
	case TypeDefinitionForeignKey:
		return "foreign_key"
	case TypeDefinitionSourceField:
		return "source_field"
	case TypeDefinitionAtribute:
		return "atribute"
	case TypeDefinitionCreatedAtField:
		return "created_at_field"
	case TypeDefinitionUpdatedAtField:
		return "updated_at_field"
	case TypeDefinitionStatusField:
		return "status_field"
	case TypeDefinitionSystemKeyField:
		return "system_key_field"
	case TypeDefinitionIndexField:
		return "index_field"
	case TypeDefinitionFullTextField:
		return "full_text_field"
	case TypeDefinitionProjectField:
		return "project_field"
	case TypeDefinitionHidden:
		return "hidden"
	case TypeDefinitionRequired:
		return "required"
	case TypeDefinitionRelation:
		return "relation"
	case TypeDefinitionRollup:
		return "rollup"
	case TypeDefinitionDetail:
		return "detail"
	case TypeDefinitionHistory:
		return "history"
	case TypeDefinitionModel:
		return "model"
	case TypeDefinitionProjectModel:
		return "project_model"
	}

	return ""
}

/**
* toTypeData
* @param val interface{}
* @return TypeData
**/
func toTypeData(val interface{}) (TypeData, error) {
	switch v := val.(type) {
	case int:
		return TypeData(v), nil
	case float64:
		i := int(v)
		return TypeData(i), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return TypeDataNone, err
		}

		return TypeData(i), nil
	}

	return TypeDataNone, mistake.New("invalid type TypeData")
}

func toArrayString(val interface{}) ([]string, error) {
	switch v := val.(type) {
	case []string:
		return v, nil
	case []interface{}:
		fields := make([]string, 0)
		for _, field := range v {
			fields = append(fields, field.(string))
		}

		return fields, nil
	case string:
		return strings.Split(v, ","), nil
	}

	return nil, mistake.New("invalid type []string")
}

func toArrayInterface(val interface{}) ([]interface{}, error) {
	switch v := val.(type) {
	case []interface{}:
		return v, nil
	case []string:
		fields := make([]interface{}, 0)
		for _, field := range v {
			fields = append(fields, field)
		}

		return fields, nil
	}

	return nil, mistake.New("invalid type []interface{}")
}
func toMapString(val interface{}) (map[string]string, error) {
	switch v := val.(type) {
	case map[string]string:
		return v, nil
	case map[string]interface{}:
		result := make(map[string]string)
		for k, v := range v {
			result[k] = v.(string)
		}
		return result, nil
	}

	return nil, mistake.New("invalid type map[string]string")
}

/**
* defineColumns
* @param tp int, args ...interface{}
**/
func (s *Model) defineColumns(tp int, args ...interface{}) error {
	tpe := TypeDefinition(tp)
	switch tpe {
	case TypeDefinitionColumn:
		tpData, err := toTypeData(args[1])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[1], err.Error())
		}

		s.defineColumn(args[0].(string), tpData)
	case TypeDefinitionIndex:
		columns, err := toArrayString(args[1])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[1], err.Error())
		}

		s.defineIndex(args[0].(bool), columns)
	case TypeDefinitionUnique:
		colums, err := toArrayString(args[0])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[0], err.Error())
		}

		s.defineUnique(colums)
	case TypeDefinitionPrimaryKey:
		colums, err := toArrayString(args[0])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[0], err.Error())
		}

		s.definePrimaryKey(colums)
	case TypeDefinitionPrimaryKeyField:
		s.definePrimaryKeyField()
	case TypeDefinitionForeignKey:
		fks, err := toMapString(args[0])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[0], err.Error())
		}

		s.defineForeignKey(fks, args[1].(string), args[2].(bool), args[3].(bool))
	case TypeDefinitionSourceField:
		s.defineSourceField()
	case TypeDefinitionAtribute:
		tpData, err := toTypeData(args[1])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[1], err.Error())
		}

		s.defineAtribute(args[0].(string), tpData)
	case TypeDefinitionCreatedAtField:
		s.defineCreatedAtField()
	case TypeDefinitionUpdatedAtField:
		s.defineUpdatedAtField()
	case TypeDefinitionStatusField:
		s.defineStatusField()
	case TypeDefinitionSystemKeyField:
		s.defineSystemKeyField()
	case TypeDefinitionIndexField:
		s.defineIndexField()
	case TypeDefinitionFullTextField:
		fields, err := toArrayString(args[1])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[1], err.Error())
		}

		s.defineFullText(args[0].(string), fields)
	case TypeDefinitionProjectField:
		s.defineProjectField()
	case TypeDefinitionHidden:
		colums, err := toArrayString(args[0])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[0], err.Error())
		}

		s.defineHidden(colums)
	case TypeDefinitionRequired:
		colums, err := toArrayString(args[0])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[0], err.Error())
		}

		s.defineRequired(colums)
	case TypeDefinitionRelation:
		fks, err := toMapString(args[2])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[2], err.Error())
		}

		s.defineRelation(args[0].(string), args[1].(string), fks, args[3].(int))
	case TypeDefinitionRollup:
		fks, err := toMapString(args[2])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[2], err.Error())
		}

		rollups, err := toArrayInterface(args[3])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[3], err.Error())
		}

		s.defineRollup(args[0].(string), args[1].(string), fks, rollups)
	case TypeDefinitionDetail:
		fks, err := toMapString(args[1])
		if err != nil {
			return mistake.Newf("invalid type: %v error: %s", args[1], err.Error())
		}

		s.defineDetail(args[0].(string), fks, args[2].(int))
	case TypeDefinitionHistory:
		s.defineHistory(args[0].(string), args[1].(int))
	case TypeDefinitionModel:
		s.defineModel()
	case TypeDefinitionProjectModel:
		s.defineProjectModel()
	}

	return nil
}

/**
* setDefine
* @param tp TypeDefinition, args ...any
**/
func (s *Model) setDefine(tp TypeDefinition, args ...any) {
	s.Definitions = append(s.Definitions, et.Json{
		"tp":   tp,
		"type": tp.Str(),
		"args": args,
	})
}

/**
* defineColumnIdx
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) defineColumnIdx(name string, typeData TypeData, idx int) *Column {
	result := s.getColumn(name)
	if result != nil {
		return result
	}

	result = newColumn(s, name, "", TpColumn, typeData, typeData.DefaultValue())
	if idx == -1 {
		s.Columns = append(s.Columns, result)
	} else {
		s.Columns = append(s.Columns[:idx], append([]*Column{result}, s.Columns[idx:]...)...)
	}

	return result
}

/**
* defineColumnIdx
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) defineColumn(name string, typeData TypeData) *Column {
	idx := s.sourceIdx()
	return s.defineColumnIdx(name, typeData, idx)
}

/**
* defineIndex
* @param sort bool, colums []string
* @return *Model
**/
func (s *Model) defineIndex(sort bool, colums []string) *Model {
	cols := s.getColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		if col.TypeColumn == TpColumn {
			idx := newIndex(col, sort)
			name := strs.Format("%s_%s_idx", s.Name, col.Name)
			s.Indices[name] = idx
		}
	}

	return s
}

/**
* defineUnique
* @param colums []string
* @return *Model
**/
func (s *Model) defineUnique(colums []string) *Model {
	cols := s.getColumns(colums...)
	if len(cols) == 0 {
		return s
	}

	for _, col := range cols {
		idx := newIndex(col, true)
		name := strs.Format("%s_%s_idx", s.Name, col.Name)
		s.Uniques[name] = idx
	}

	return s
}

/**
* definePrimaryKey
* @param name string
* @return *Column
**/
func (s *Model) definePrimaryKey(primaryKeys []string) *Model {
	for _, primaryKey := range primaryKeys {
		col := s.getColumn(primaryKey)
		if col != nil {
			s.Required[col.Name] = true
			s.PrimaryKeys[col.Name] = col
		}
	}

	return s
}

/**
* definePrimaryKeyField
* @return *Column
**/
func (s *Model) definePrimaryKeyField() *Column {
	result := s.defineColumn(PRIMARYKEY, PrimaryKeyField.TypeData())
	s.definePrimaryKey([]string{PRIMARYKEY})

	return result
}

/**
* defineForeignKey
* @param fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool
* @return *Relation
**/
func (s *Model) defineForeignKey(fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool) *Relation {
	with := s.Db.GetModel(withName)
	if with == nil {
		return nil
	}

	result := &Relation{
		With:            with,
		Fk:              make(map[string]string),
		Limit:           0,
		OnDeleteCascade: onDeleteCascade,
		OnUpdateCascade: onUpdateCascade,
	}

	for fkn, pkn := range fks {
		pk := with.getColumn(pkn)
		if pk == nil {
			continue
		}

		result.Fk[fkn] = pk.Name
		fk := s.defineColumn(fkn, pk.TypeData)
		fk.Detail = result
		name := strs.Format("%s_%s_fk", s.Name, with.Name)
		s.ForeignKeys[name] = fk
	}

	return result
}

/**
* defineSourceField
* @return *Column
**/
func (s *Model) defineSourceField() *Column {
	if s.SourceField != nil {
		return s.SourceField
	}

	result := s.defineColumn(SOURCE, SourceField.TypeData())
	s.defineIndex(true, []string{SOURCE})
	s.SourceField = result

	return result
}

/**
* defineAtribute
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) defineAtribute(name string, typeData TypeData) *Column {
	s.defineSourceField()
	result := newAtribute(s, name, typeData)
	s.Columns = append(s.Columns, result)

	return result
}

/**
* defineCreatedAtField
* @return *Column
**/
func (s *Model) defineCreatedAtField() *Column {
	result := s.defineColumn(string(CreatedAtField), CreatedAtField.TypeData())
	s.defineIndex(true, []string{CreatedAtField.Str()})

	return result
}

/**
* defineUpdatedAtField
* @return *Column
**/
func (s *Model) defineUpdatedAtField() *Column {
	result := s.defineColumn(string(UpdatedAtField), UpdatedAtField.TypeData())
	s.defineIndex(true, []string{UpdatedAtField.Str()})

	return result
}

/**
* defineStatusField
* @return *Column
**/
func (s *Model) defineStatusField() *Column {
	result := s.defineColumn(string(StatusField), StatusField.TypeData())
	s.defineIndex(true, []string{StatusField.Str()})
	s.StatusField = result

	return result
}

/**
* defineSystemKeyField
* @return *Column
**/
func (s *Model) defineSystemKeyField() *Column {
	result := s.defineColumn(string(SystemKeyField), SystemKeyField.TypeData())
	s.defineIndex(true, []string{SystemKeyField.Str()})
	s.SystemKeyField = result

	return result
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) defineIndexField() *Column {
	result := s.defineColumn(string(IndexField), IndexField.TypeData())
	s.defineIndex(true, []string{IndexField.Str()})
	s.IndexField = result

	return result
}

/**
* DefineFullText
* @param fields []string
* @return language string
* @return *Column
**/
func (s *Model) defineFullText(language string, fields []string) *Column {
	cols := s.getColumnsArray(fields...)
	result := s.defineColumn(string(FullTextField), FullTextField.TypeData())
	result.FullText = &FullText{
		Language: language,
		Columns:  cols,
	}
	s.defineIndex(true, []string{FullTextField.Str()})
	s.FullTextField = result

	return result
}

/**
* defineProjectField
* @return *Column
**/
func (s *Model) defineProjectField() *Column {
	idx := -1
	for _, pk := range s.PrimaryKeys {
		min := slices.IndexFunc(s.Columns, func(e *Column) bool { return e.Name == pk.Name })
		if min != -1 && idx == -1 {
			idx = min
		} else if min != -1 && min < idx {
			idx = min
		}
	}

	result := s.defineColumnIdx(string(ProjectField), ProjectField.TypeData(), idx)
	s.defineIndex(true, []string{ProjectField.Str()})
	s.ProjectField = result

	return result
}

/**
* defineHidden
* @param colums []string
* @return *Model
**/
func (s *Model) defineHidden(colums []string) *Model {
	for _, name := range colums {
		col := s.getColumn(name)
		if col != nil {
			col.Hidden = true
		}
	}

	return s
}

/**
* defineRequired
* @param colums []string
* @return *Model
**/
func (s *Model) defineRequired(colums []string) *Model {
	for _, name := range colums {
		col := s.getColumn(name)
		if col != nil {
			s.Required[name] = true
		}
	}

	return s
}

/**
* defineRelation
* @param name, relatedTo string, fks map[string]string, limit int
* @return *Relation
**/
func (s *Model) defineRelation(name, relatedTo string, fks map[string]string, limit int) *Relation {
	with := NewModel(s.Schema, relatedTo, 1)
	with.defineForeignKey(fks, s.Name, true, true)

	result := &Relation{
		With:  with,
		Fk:    make(map[string]string),
		Limit: limit,
	}
	for fkn, pkn := range fks {
		fk := with.getColumn(fkn)
		if fk == nil {
			continue
		}

		result.Fk[pkn] = fk.Name
	}
	col := newColumn(s, name, "", TpRelatedTo, TypeDataNone, TypeDataNone.DefaultValue())
	col.Detail = result
	s.RelationsTo[with.Name] = result
	s.Columns = append(s.Columns, col)

	return result
}

/**
* defineRollup
* @param name, rollupFrom string, fks map[string]string, fields []interface{}
* @return *Rollup
**/
func (s *Model) defineRollup(name, rollupFrom string, fks map[string]string, fields []interface{}) *Rollup {
	source := s.Db.GetModel(rollupFrom)
	if source == nil {
		return nil
	}

	result := &Rollup{
		Source: source,
		Fk:     make(map[string]string),
		Fields: fields,
	}

	for fkn, pkn := range fks {
		pk := source.getColumn(pkn)
		if pk == nil {
			return nil
		}

		result.Fk[fkn] = pk.Name
	}

	col := newColumn(s, name, "", TpRollup, TypeDataNone, TypeDataNone.DefaultValue())
	col.Rollup = result
	s.Columns = append(s.Columns, col)
	s.Rollups[name] = result

	return result
}

/**
* defineDetail
* @param name string, fks map[string]string, limit int
* @return *Model
**/
func (s *Model) defineDetail(name string, fks map[string]string, limit int) *Model {
	relatedTo := s.Name + "_" + name
	result := s.defineRelation(name, relatedTo, fks, limit)
	s.Details[name] = result

	return result.With
}

/**
* defineHistory
* @param pkn string, limit int
* @return *Model
**/
func (s *Model) defineHistory(pkn string, limit int) *Model {
	relatedTo := s.Name + "_" + HISTORYCAL
	result := s.defineRelation(HISTORYCAL, relatedTo, map[string]string{"history_id": pkn}, limit)
	with := result.With
	with.defineColumn(CREATED_AT, CreatedAtField.TypeData())
	with.defineColumn(SYSID, SystemKeyField.TypeData())
	with.defineColumn(HISTORYCAL, SourceField.TypeData())
	with.defineColumn(INDEX, IndexField.TypeData())
	with.defineIndex(true, []string{SYSID, INDEX})
	s.History = result

	return with
}

/**
* defineModel
* @return *Model
**/
func (s *Model) defineModel() *Model {
	s.defineCreatedAtField()
	s.defineUpdatedAtField()
	s.defineStatusField()
	s.definePrimaryKeyField()
	s.defineSourceField()
	s.defineIndexField()
	s.defineSystemKeyField()

	return s
}

/**
* defineProjectModel
* @return *Model
**/
func (s *Model) defineProjectModel() *Model {
	s.defineCreatedAtField()
	s.defineUpdatedAtField()
	s.defineProjectField()
	s.defineStatusField()
	s.definePrimaryKeyField()
	s.defineSourceField()
	s.defineIndexField()
	s.defineSystemKeyField()

	return s
}

/**
* DefineIntegrity
**/
func (s *Model) DefineIntegrity() {
	s.Integrity = true
}

/**
* DefineColumn
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	s.setDefine(TypeDefinitionColumn, name, typeData)
	return s.defineColumn(name, typeData)
}

/**
* DefineIndex
* @param sort bool, colums []string
* @return *Model
**/
func (s *Model) DefineIndex(sort bool, colums ...string) *Model {
	s.setDefine(TypeDefinitionIndex, sort, colums)
	return s.defineIndex(sort, colums)
}

/**
* DefineUnique
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineUnique(colums ...string) *Model {
	s.setDefine(TypeDefinitionUnique, colums)
	return s.defineUnique(colums)
}

/**
* DefinePrimaryKey
* @param colums ...string
* @return *Model
**/
func (s *Model) DefinePrimaryKey(colums ...string) *Model {
	s.setDefine(TypeDefinitionPrimaryKey, colums)
	return s.definePrimaryKey(colums)
}

/**
* DefinePrimaryKeyField
* @return *Column
**/
func (s *Model) DefinePrimaryKeyField() *Column {
	s.setDefine(TypeDefinitionPrimaryKeyField)
	return s.definePrimaryKeyField()
}

/**
* DefineForeignKey
* @param fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool
* @return *Relation
**/
func (s *Model) DefineForeignKey(fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool) *Relation {
	s.setDefine(TypeDefinitionForeignKey, fks, withName, onDeleteCascade, onUpdateCascade)
	return s.defineForeignKey(fks, withName, onDeleteCascade, onUpdateCascade)
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField() *Column {
	s.setDefine(TypeDefinitionSourceField)
	return s.defineSourceField()
}

/**
* DefineAtribute
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineAtribute(name string, typeData TypeData) *Column {
	s.setDefine(TypeDefinitionAtribute, name, typeData)
	return s.defineAtribute(name, typeData)
}

/**
* DefineCreatedAtField
* @return *Column
**/
func (s *Model) DefineCreatedAtField() *Column {
	s.setDefine(TypeDefinitionCreatedAtField)
	return s.defineCreatedAtField()
}

/**
* DefineUpdatedAtField
* @return *Column
**/
func (s *Model) DefineUpdatedAtField() *Column {
	s.setDefine(TypeDefinitionUpdatedAtField)
	return s.defineUpdatedAtField()
}

/**
* DefineStatusField
* @return *Column
**/
func (s *Model) DefineStatusField() *Column {
	s.setDefine(TypeDefinitionStatusField)
	return s.defineStatusField()
}

/**
* DefineSystemKeyField
* @return *Column
**/
func (s *Model) DefineSystemKeyField() *Column {
	s.setDefine(TypeDefinitionSystemKeyField)
	return s.defineSystemKeyField()
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) DefineIndexField() *Column {
	s.setDefine(TypeDefinitionIndexField)
	return s.defineIndexField()
}

/**
* DefineFullText
* @param language string, fields []string
* @return *Column
**/
func (s *Model) DefineFullText(language string, fields []string) *Column {
	s.setDefine(TypeDefinitionFullTextField, language, fields)
	return s.defineFullText(language, fields)
}

/**
* DefineProjectField
* @return *Column
**/
func (s *Model) DefineProjectField() *Column {
	s.setDefine(TypeDefinitionProjectField)
	return s.defineProjectField()
}

/**
* DefineHidden
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineHidden(colums ...string) *Model {
	s.setDefine(TypeDefinitionHidden, colums)
	return s.defineHidden(colums)
}

/**
* DefineRequired
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineRequired(colums ...string) *Model {
	s.setDefine(TypeDefinitionRequired, colums)
	return s.defineRequired(colums)
}

/**
* DefineRelation
* @param name, relatedTo string, fks map[string]string, limit int
* @return *Relation
**/
func (s *Model) DefineRelation(name, relatedTo string, fks map[string]string, limit int) *Relation {
	s.setDefine(TypeDefinitionRelation, name, relatedTo, fks, limit)
	return s.defineRelation(name, relatedTo, fks, limit)
}

/**
* DefineRollup
* @param name, rollupFrom string, fks map[string]string, properties []interface{}
* @return *Rollup
**/
func (s *Model) DefineRollup(name, rollupFrom string, fks map[string]string, properties []interface{}) *Rollup {
	s.setDefine(TypeDefinitionRollup, name, rollupFrom, fks, properties)
	return s.defineRollup(name, rollupFrom, fks, properties)
}

/**
* DefineDetail
* @param name string, fks map[string]string, limit int
* @return *Model
**/
func (s *Model) DefineDetail(name string, fks map[string]string, limit int) *Model {
	s.setDefine(TypeDefinitionDetail, name, fks, limit)
	return s.defineDetail(name, fks, limit)
}

/**
* DefineHistory
* @param pkn string, limit int
* @return *Model
**/
func (s *Model) DefineHistory(pkn string, limit int) *Model {
	s.setDefine(TypeDefinitionHistory, pkn, limit)
	return s.defineHistory(pkn, limit)
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	s.setDefine(TypeDefinitionModel)
	return s.defineModel()
}

/**
* DefineProjectModel
* @return *Model
**/
func (s *Model) DefineProjectModel() *Model {
	s.setDefine(TypeDefinitionProjectModel)
	return s.defineProjectModel()
}

/**
* DefineGenerated
* @param name string, fn GeneratedFunction
* @return *Column
**/
func (s *Model) DefineGenerated(name string, fn GeneratedFunction) *Column {
	result := newColumn(s, name, "", TpGenerated, TypeDataNone, TypeDataNone.DefaultValue())
	result.GeneratedFunction = fn
	s.Columns = append(s.Columns, result)
	s.GeneratedFields = append(s.GeneratedFields, result)

	return result
}

/**
* DefineEvent
* @param tp TypeEvent, event Event
**/
func (s *Model) DefineEvent(tp TypeEvent, event Event) {
	switch tp {
	case EventInsert:
		s.eventsInsert = append(s.eventsInsert, event)
	case EventUpdate:
		s.eventsUpdate = append(s.eventsUpdate, event)
	case EventDelete:
		s.eventsDelete = append(s.eventsDelete, event)
	}
}

/**
* DefineEventError
* @param event EventError
**/
func (s *Model) DefineEventError(event EventError) {
	s.eventError = append(s.eventError, event)
}
