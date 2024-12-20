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
	result.DefineColumn(UpdatedAtField.Str(), UpdatedAtField.TypeData())
	result.DefineColumn(ProjectField.Str(), ProjectField.TypeData())
	result.DefineColumn(StateField.Str(), StateField.TypeData())
	result.DefineColumn(KeyField.Str(), KeyField.TypeData())
	result.DefineColumn(ClassField.Str(), ClassField.TypeData())
	result.DefineColumn(SourceField.Str(), SourceField.TypeData())
	result.DefineColumn(SystemKeyField.Str(), SystemKeyField.TypeData())
	result.DefineColumn(IndexField.Str(), IndexField.TypeData())
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
