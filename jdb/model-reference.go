package jdb

type TypeRelation int

const (
	RelationOneToOne TypeRelation = iota
	RelationOneToMany
	RelationManyToMany
	RelationManyToOne
)

type Reference struct {
	Pk              []*Column
	Fk              []*Column
	Detail          *Model
	TypeRelation    TypeRelation
	OnDeleteCascade bool
	OnUpdateCascade bool
}

func NewReference(from *Model, to *Model, tp TypeRelation, pks []*Column, fks []*Column) *Reference {
	result := &Reference{
		Pk:           pks,
		Fk:           fks,
		Detail:       to,
		TypeRelation: tp,
	}
	from.References = append(from.References, result)

	return result
}
