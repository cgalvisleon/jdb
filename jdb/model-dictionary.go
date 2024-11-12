package jdb

type Dictionary struct {
	Model   *Model
	Key     string
	Value   interface{}
	Columns []*Column
}

/**
* NewDictionary
* @param key string
* @param value interface{}
* @return *Dictionary
**/
func NewDictionary(model *Model, key string, value interface{}) *Dictionary {
	return &Dictionary{
		Model:   model,
		Key:     key,
		Value:   value,
		Columns: []*Column{},
	}
}

/**
* DefineAtribute
* @param name string
* @param typeData TypeData
* @param def interface{}
* @return *Model
**/
func (s *Dictionary) DefineAtribute(name string, typeData TypeData, def interface{}) *Column {
	col := newColumn(s.Model, name, "", TpAtribute, typeData, def)
	s.Columns = append(s.Columns, col)

	return col
}
