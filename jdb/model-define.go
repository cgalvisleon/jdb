package jdb

import (
	"fmt"
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
	TypeDefinitionSource
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
	TypeDefinitionValues
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
	case TypeDefinitionSource:
		return "source"
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
	case TypeDefinitionValues:
		return "values"
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

	return TypeDataNone, mistake.Newf("invalid type: %T to TypeData", val)
}

/**
* toArrayString
* @param val interface{}
* @return []string, error
**/
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

	return nil, mistake.Newf("invalid type: %T to []string", val)
}

/**
* toMapString
* @param val interface{}
* @return map[string]string, error
**/
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

	return nil, mistake.Newf("invalid type: %T to map[string]string", val)
}

/**
* toMapInt
* @param val interface{}
* @return int, error
**/
func toMapInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		n, err := strconv.Atoi(v)
		return n, err
	}

	return 0, fmt.Errorf("invalid type: %T to int", val)
}

/**
* toShowRollup
* @param val interface{}
* @return ShowRollup, error
**/
func toShowRollup(val interface{}) (ShowRollup, error) {
	switch v := val.(type) {
	case int:
		if v == 0 {
			return ShowAtrib, nil
		}
		return ShowObject, nil
	case float64:
		i := int(v)
		if i == 0 {
			return ShowAtrib, nil
		}
		return ShowObject, nil
	case string:
		switch v {
		case "object":
			return ShowObject, nil
		default:
			return ShowAtrib, nil
		}
	default:
		return ShowAtrib, nil
	}
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
			return err
		}

		s.defineColumn(args[0].(string), tpData)
	case TypeDefinitionIndex:
		columns, err := toArrayString(args[1])
		if err != nil {
			return err
		}

		s.defineIndex(args[0].(bool), columns)
	case TypeDefinitionUnique:
		colums, err := toArrayString(args[0])
		if err != nil {
			return err
		}

		s.defineUnique(colums)
	case TypeDefinitionPrimaryKey:
		colums, err := toArrayString(args[0])
		if err != nil {
			return err
		}

		s.definePrimaryKey(colums)
	case TypeDefinitionPrimaryKeyField:
		s.definePrimaryKeyField()
	case TypeDefinitionForeignKey:
		fks, err := toMapString(args[0])
		if err != nil {
			return err
		}

		s.defineForeignKey(fks, args[1].(string), args[2].(bool), args[3].(bool))
	case TypeDefinitionSource:
		s.defineSource(args[0].(string))
	case TypeDefinitionSourceField:
		s.defineSourceField()
	case TypeDefinitionAtribute:
		tpData, err := toTypeData(args[1])
		if err != nil {
			return err
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
			return err
		}

		s.defineFullText(args[0].(string), fields)
	case TypeDefinitionProjectField:
		s.defineProjectField()
	case TypeDefinitionHidden:
		colums, err := toArrayString(args[0])
		if err != nil {
			return err
		}

		s.defineHidden(colums)
	case TypeDefinitionRequired:
		colums, err := toArrayString(args[0])
		if err != nil {
			return err
		}

		s.defineRequired(colums)
	case TypeDefinitionRelation:
		fks, err := toMapString(args[2])
		if err != nil {
			return err
		}

		limit, err := toMapInt(args[3])
		if err != nil {
			return err
		}

		s.defineRelation(args[0].(string), args[1].(string), fks, limit)
	case TypeDefinitionRollup:
		fks, err := toMapString(args[2])
		if err != nil {
			return err
		}

		fields, err := toArrayString(args[3])
		if err != nil {
			return err
		}

		showRollup, err := toShowRollup(args[4])
		if err != nil {
			return err
		}

		s.defineRollup(args[0].(string), args[1].(string), fks, fields, showRollup)
	case TypeDefinitionDetail:
		fks, err := toMapString(args[1])
		if err != nil {
			return err
		}

		limit, err := toMapInt(args[2])
		if err != nil {
			return err
		}

		s.defineDetail(args[0].(string), fks, limit)
	case TypeDefinitionValues:
		s.defineValues(args[0].(string), args[1])
	case TypeDefinitionModel:
		s.defineModel()
	case TypeDefinitionProjectModel:
		s.defineProjectModel()
	}

	return nil
}

/**
* setDefine
* @param name string, tp TypeDefinition, args ...any
**/
func (s *Model) setDefine(name string, tp TypeDefinition, args ...any) {
	s.Definitions[name] = et.Json{
		"tp":   tp,
		"type": tp.Str(),
		"args": args,
	}
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
		s.addColumn(result)
	} else {
		s.addColumnToIdx(result, idx)
	}

	return result
}

/**
* defineColumn
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
* @param primaryKeys []string
* @return *Model
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
	result.IsKeyfield = true
	s.definePrimaryKey([]string{PRIMARYKEY})

	return result
}

/**
* defineForeignKey
* @param fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool
* @return *Relation
**/
func (s *Model) defineForeignKey(fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool) *Relation {
	with := s.GetModel(withName)
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

	from := &Relation{
		With:            s,
		Fk:              make(map[string]string),
		Limit:           0,
		OnDeleteCascade: onDeleteCascade,
		OnUpdateCascade: onUpdateCascade,
	}

	name := strs.Format("%s_%s_fk", s.Name, with.Name)
	for fkn, pkn := range fks {
		pk := with.getColumn(pkn)
		if pk == nil {
			continue
		}

		result.Fk[fkn] = pk.Name
		from.Fk[pk.Name] = fkn
		fk := s.defineColumn(fkn, pk.TypeData)
		fk.Detail = result
		s.Required[fkn] = true
		s.defineIndex(true, []string{fkn})
	}
	s.ForeignKeys[name] = result
	with.RelationsFrom[s.Name] = from

	return result
}

/**
* defineSource
* @param name string
* @return *Column
**/
func (s *Model) defineSource(name string) *Column {
	result := s.defineColumn(name, SourceField.TypeData())
	s.defineIndex(true, []string{name})
	s.SourceField = result

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

	return s.defineSource(SourceField.Str())
}

/**
* defineAtribute
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) defineAtribute(name string, typeData TypeData) *Column {
	s.defineSourceField()
	result := newAtribute(s, name, typeData)
	s.addColumn(result)

	return result
}

/**
* defineCreatedAtField
* @return *Column
**/
func (s *Model) defineCreatedAtField() *Column {
	result := s.defineColumn(string(CreatedAtField), CreatedAtField.TypeData())
	s.defineIndex(true, []string{CreatedAtField.Str()})
	s.CreatedAtField = result

	return result
}

/**
* defineUpdatedAtField
* @return *Column
**/
func (s *Model) defineUpdatedAtField() *Column {
	result := s.defineColumn(string(UpdatedAtField), UpdatedAtField.TypeData())
	s.defineIndex(true, []string{UpdatedAtField.Str()})
	s.UpdatedAtField = result

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
	with := NewModel(s.schema, relatedTo, 1)
	with.defineForeignKey(fks, s.Name, true, true)

	col := newColumn(s, name, "", TpRelatedTo, TypeDataNone, TypeDataNone.DefaultValue())
	col.Detail = &Relation{
		With:  with,
		Fk:    fks,
		Limit: limit,
	}
	s.RelationsTo[name] = col.Detail
	s.addColumn(col)

	return col.Detail
}

/**
* defineRollup
* @param name, rollupFrom string, fks map[string]string, fields []string, showRollup ShowRollup
* @return *Column
**/
func (s *Model) defineRollup(name, rollupFrom string, fks map[string]string, fields []string, showRollup ShowRollup) *Column {
	with := NewModel(s.schema, rollupFrom, 1)
	result := &Rollup{
		With:   with,
		Fk:     fks,
		Fields: fields,
		Show:   showRollup,
	}

	for fkn, pkn := range fks {
		fk := s.getColumn(fkn)
		if fk != nil {
			continue
		}

		pk := with.getColumn(pkn)
		if pk != nil {
			fk = s.defineColumn(fkn, pk.TypeData)
			s.defineIndex(true, []string{fkn})
			s.addColumn(fk)
		}
	}

	col := newColumn(s, name, "", TpRollup, TypeDataNone, TypeDataNone.DefaultValue())
	col.Rollup = result
	s.addColumn(col)

	return col
}

/**
* defineDetail
* @param name string, fks map[string]string, limit int
* @return *Model
**/
func (s *Model) defineDetail(name string, fks map[string]string, limit int) *Model {
	relatedTo := s.Name + "_" + name
	result := s.defineRelation(name, relatedTo, fks, limit)

	return result.With
}

/**
* defineMultiSelect
* @param name string, fks map[string]string
* @return *Model
**/
func (s *Model) defineMultiSelect(name string, fks map[string]string) *Model {
	relatedTo := s.Name + "_" + name
	result := s.defineRelation(name, relatedTo, fks, 30)
	result.With.definePrimaryKeyField()
	result.With.defineColumn(CHECKED, TypeDataCheckbox)
	result.With.defineCreatedAtField()
	result.With.defineSourceField()
	result.With.defineSystemKeyField()
	result.With.defineIndexField()
	primaryKeys := []string{}
	for fkn := range fks {
		primaryKeys = append(primaryKeys, fkn)
	}
	primaryKeys = append(primaryKeys, KEY)
	result.With.definePrimaryKey(primaryKeys)
	result.IsMultiSelect = true

	return result.With
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
	s.defineSystemKeyField()
	s.defineIndexField()

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
	s.defineSystemKeyField()
	s.defineIndexField()

	return s
}

/**
* defineValues
* @param name string, values interface{}
* @return *Column
**/
func (s *Model) defineValues(name string, values interface{}) *Column {
	col := s.defineAtribute(name, TypeDataObject)
	col.Values = values

	return col
}

/**
* DefineIntegrity
**/
func (s *Model) DefineIntegrity() *Model {
	s.Integrity = true
	return s
}

/**
* DefineColumn
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineColumn(name string, typeData TypeData) *Column {
	key := fmt.Sprintf("column_%v", name)
	s.setDefine(key, TypeDefinitionColumn, name, typeData)
	return s.defineColumn(name, typeData)
}

/**
* DefineIndex
* @param sort bool, colums []string
* @return *Model
**/
func (s *Model) DefineIndex(sort bool, colums ...string) *Model {
	key := fmt.Sprintf("index_%v", sort)
	s.setDefine(key, TypeDefinitionIndex, sort, colums)
	return s.defineIndex(sort, colums)
}

/**
* DefineUnique
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineUnique(colums ...string) *Model {
	key := "unique"
	s.setDefine(key, TypeDefinitionUnique, colums)
	return s.defineUnique(colums)
}

/**
* DefinePrimaryKey
* @param colums ...string
* @return *Model
**/
func (s *Model) DefinePrimaryKey(colums ...string) *Model {
	key := "primary_key"
	s.setDefine(key, TypeDefinitionPrimaryKey, colums)
	return s.definePrimaryKey(colums)
}

/**
* DefinePrimaryKeyField
* @return *Model
**/
func (s *Model) DefinePrimaryKeyField() *Model {
	key := "primary_key_field"
	s.setDefine(key, TypeDefinitionPrimaryKeyField)
	s.definePrimaryKeyField()
	return s
}

/**
* DefineForeignKey
* @param fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool
* @return *Model
**/
func (s *Model) DefineForeignKey(fks map[string]string, withName string, onDeleteCascade, onUpdateCascade bool) *Model {
	key := fmt.Sprintf("foreign_key_%v", withName)
	s.setDefine(key, TypeDefinitionForeignKey, fks, withName, onDeleteCascade, onUpdateCascade)
	s.defineForeignKey(fks, withName, onDeleteCascade, onUpdateCascade)
	return s
}

/**
* DefineSource
* @param name string
* @return *Column
**/
func (s *Model) DefineSource(name string) *Column {
	key := fmt.Sprintf("source_%v", name)
	s.setDefine(key, TypeDefinitionSource, name)
	return s.defineSource(name)
}

/**
* DefineSourceField
* @return *Column
**/
func (s *Model) DefineSourceField() *Column {
	key := "source_field"
	s.setDefine(key, TypeDefinitionSourceField)
	return s.defineSourceField()
}

/**
* DefineAtribute
* @param name string, typeData TypeData
* @return *Column
**/
func (s *Model) DefineAtribute(name string, typeData TypeData) *Column {
	key := fmt.Sprintf("atribute_%v", name)
	s.setDefine(key, TypeDefinitionAtribute, name, typeData)
	return s.defineAtribute(name, typeData)
}

/**
* DefineCreatedAtField
* @return *Column
**/
func (s *Model) DefineCreatedAtField() *Column {
	key := "created_at_field"
	s.setDefine(key, TypeDefinitionCreatedAtField)
	return s.defineCreatedAtField()
}

/**
* DefineUpdatedAtField
* @return *Column
**/
func (s *Model) DefineUpdatedAtField() *Column {
	key := "updated_at_field"
	s.setDefine(key, TypeDefinitionUpdatedAtField)
	return s.defineUpdatedAtField()
}

/**
* DefineStatusField
* @return *Column
**/
func (s *Model) DefineStatusField() *Column {
	key := "status_field"
	s.setDefine(key, TypeDefinitionStatusField)
	return s.defineStatusField()
}

/**
* DefineSystemKeyField
* @return *Column
**/
func (s *Model) DefineSystemKeyField() *Column {
	key := "system_key_field"
	s.setDefine(key, TypeDefinitionSystemKeyField)
	return s.defineSystemKeyField()
}

/**
* DefineIndexField
* @return *Column
**/
func (s *Model) DefineIndexField() *Column {
	key := "index_field"
	s.setDefine(key, TypeDefinitionIndexField)
	return s.defineIndexField()
}

/**
* DefineFullText
* @param language string, fields []string
* @return *Model
**/
func (s *Model) DefineFullText(language string, fields []string) *Column {
	key := fmt.Sprintf("full_text_%v", language)
	s.setDefine(key, TypeDefinitionFullTextField, language, fields)
	return s.defineFullText(language, fields)
}

/**
* DefineProjectField
* @return *Column
**/
func (s *Model) DefineProjectField() *Column {
	key := "project_field"
	s.setDefine(key, TypeDefinitionProjectField)
	return s.defineProjectField()
}

/**
* DefineHidden
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineHidden(colums ...string) *Model {
	key := "hidden"
	s.setDefine(key, TypeDefinitionHidden, colums)
	s.defineHidden(colums)
	return s
}

/**
* DefineRequired
* @param colums ...string
* @return *Model
**/
func (s *Model) DefineRequired(colums ...string) *Model {
	key := "required"
	s.setDefine(key, TypeDefinitionRequired, colums)
	s.defineRequired(colums)
	return s
}

/**
* DefineRelation
* @param name, relatedTo string, fks map[string]string, limit int
* @return *Model
**/
func (s *Model) DefineRelation(name, relatedTo string, fks map[string]string, limit int) *Relation {
	key := fmt.Sprintf("relation_%v", name)
	s.setDefine(key, TypeDefinitionRelation, name, relatedTo, fks, limit)
	return s.defineRelation(name, relatedTo, fks, limit)
}

/**
* DefineRollup
* @param name, rollupFrom string, fks map[string]string, field string
* @return *Model
**/
func (s *Model) DefineRollup(name, rollupFrom string, fks map[string]string, field string) *Column {
	key := fmt.Sprintf("rollup_%v", name)
	s.setDefine(key, TypeDefinitionRollup, name, rollupFrom, fks, field, ShowAtrib)
	return s.defineRollup(name, rollupFrom, fks, []string{field}, ShowAtrib)
}

/**
* DefineObject
* @param name, rollupFrom string, fks map[string]string, field string
* @return *Model
**/
func (s *Model) DefineObject(name, rollupFrom string, fks map[string]string, fields []string) *Model {
	key := fmt.Sprintf("object_%v", name)
	s.setDefine(key, TypeDefinitionRollup, name, rollupFrom, fks, fields, ShowObject)
	s.defineRollup(name, rollupFrom, fks, fields, ShowObject)
	return s
}

/**
* DefineDetail
* @param name string, fks map[string]string, limit int
* @return *Model
**/
func (s *Model) DefineDetail(name string, fks map[string]string, limit int) *Model {
	key := fmt.Sprintf("detail_%v", name)
	s.setDefine(key, TypeDefinitionDetail, name, fks, limit)
	return s.defineDetail(name, fks, limit)
}

/**
* DefineMultiSelect
* @param name string, fks map[string]string
* @return *Model
**/
func (s *Model) DefineMultiSelect(name string, fks map[string]string) *Model {
	key := fmt.Sprintf("multi_select_%v", name)
	s.setDefine(key, TypeDefinitionDetail, name, fks)
	return s.defineMultiSelect(name, fks)
}

/**
* DefineModel
* @return *Model
**/
func (s *Model) DefineModel() *Model {
	key := "model"
	s.setDefine(key, TypeDefinitionModel)
	return s.defineModel()
}

/**
* DefineProjectModel
* @return *Model
**/
func (s *Model) DefineProjectModel() *Model {
	key := "project_model"
	s.setDefine(key, TypeDefinitionProjectModel)
	return s.defineProjectModel()
}

/**
* DefineCalc
* @param name string, fn DataFunction
* @return Model
**/
func (s *Model) DefineCalc(name string, fn DataFunction) *Model {
	result := s.getColumn(name)
	if result != nil {
		result.CalcFunction = fn
		return s
	}

	result = newColumn(s, name, "", TpCalc, TypeDataNone, TypeDataNone.DefaultValue())
	result.CalcFunction = fn
	s.addColumn(result)

	return s
}

/**
* DefineEvent
* @param tp TypeEvent, event Event
* @return Model
**/
func (s *Model) DefineEvent(tp TypeEvent, event Event) *Model {
	switch tp {
	case EventInsert:
		s.eventsInsert = append(s.eventsInsert, event)
	case EventUpdate:
		s.eventsUpdate = append(s.eventsUpdate, event)
	case EventDelete:
		s.eventsDelete = append(s.eventsDelete, event)
	}

	return s
}

/**
* DefineFunc
* @param tp TypeEvent, jsCode string
* @return Model
**/
func (s *Model) DefineFunc(tp TypeEvent, jsCode string) *Model {
	switch tp {
	case EventInsert:
		s.FuncInsert = append(s.FuncInsert, jsCode)
	case EventUpdate:
		s.FuncUpdate = append(s.FuncUpdate, jsCode)
	case EventDelete:
		s.FuncDelete = append(s.FuncDelete, jsCode)
	}

	return s
}

/**
* DefineEventError
* @param event EventError
* @return Model
**/
func (s *Model) DefineEventError(event EventError) *Model {
	s.eventError = append(s.eventError, event)
	return s
}

/**
* DefineValues
* @param name string, values interface{}
* @return *Column
**/
func (s *Model) DefineValues(name string, values interface{}) *Column {
	key := fmt.Sprintf("values_%v", name)
	s.setDefine(key, TypeDefinitionValues, name, values)
	return s.defineValues(name, values)
}

/**
* defineFields
* @param fields et.Json
**/
func (s *Model) defineFields(fields et.Json) {
	for key := range fields {
		definition := fields.Json(key)
		var typeColumn TypeColumn
		var typeData TypeData
		if definition["kind"] == nil {
			tipe := definition.ValStr("text", "type")
			typeColumn, typeData = StrToKindType(tipe)
		} else {
			kind := definition.Str("kind")
			typeColumn = StrsToTypeColumn(kind)
			typeData = StrsToTypeData(definition.ValStr("text", "type"))
		}
		var field *Column
		switch typeColumn {
		case TpColumn:
			field = s.DefineColumn(key, typeData)
		case TpAtribute:
			field = s.DefineAtribute(key, typeData)
		case TpRelatedTo:
			limit := definition.ValInt(30, "limit")
			relatedTo := definition.Str("related_to")
			fks := map[string]string{}
			for key, value := range definition.Json("foreign_keys") {
				fks[key] = fmt.Sprintf("%v", value)
			}
			relation := s.DefineRelation(key, relatedTo, fks, limit)
			fields := definition.Json("fields")
			relation.With.setFields(fields)
		case TpRollup:
			rollupFrom := definition.Str("rollup_from")
			fieldName := definition.Str("field")
			fks := map[string]string{}
			for key, value := range definition.Json("foreign_keys") {
				fks[key] = fmt.Sprintf("%v", value)
			}
			field = s.DefineRollup(key, rollupFrom, fks, fieldName)
		}

		if field != nil {
			field.Description = definition.ValStr("", "description")
			field.Default = definition["default"]
			field.Max = definition.Num("max")
			field.Min = definition.Num("min")
			field.SetValue(definition["values"])
			if definition.Bool("hidden") {
				s.DefineHidden(key)
			}
			if definition.Bool("required") {
				s.DefineRequired(key)
			}
			if definition.Bool("unique") {
				s.DefineUnique(key)
			}

		}
	}
}
