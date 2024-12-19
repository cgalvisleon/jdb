package jdb

/**
* NewCollection
* @param schema *Schema
* @param name string
* @return *Model
**/
func NewCollection(schema *Schema, name string) *Model {
	result := NewModel(schema, name)
	result.DefineColumn(CreatedAtField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(UpdatedAtField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(ProjectField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(StateField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(KeyField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(ClassField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(SourceField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(SystemKeyField.Str(), CreatedAtField.TypeData())
	result.DefineColumn(IndexField.Str(), CreatedAtField.TypeData())
	result.DefineKey(KeyField.Str())
	result.DefineIndex(true,
		CreatedAtField.Str(),
		UpdatedAtField.Str(),
		ProjectField.Str(),
		StateField.Str(),
		KeyField.Str(),
		ClassField.Str(),
		SourceField.Str(),
		SystemKeyField.Str(),
		IndexField.Str(),
	)

	return result
}
