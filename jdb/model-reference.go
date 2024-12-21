package jdb

type TypeRelation int

const (
	RelationOneToOne TypeRelation = iota
	RelationOneToMany
	RelationManyToMany
	RelationManyToOne
)

type Reference struct {
	Key             *Column
	TypeRelation    TypeRelation
	To              *Column
	OnDeleteCascade bool
	OnUpdateCascade bool
}

func NewReference(key *Column, tp TypeRelation, to *Column) *Reference {
	result := &Reference{
		Key:          key,
		TypeRelation: tp,
		To:           to,
	}
	key.Model.References = append(key.Model.References, result)

	return result
}
