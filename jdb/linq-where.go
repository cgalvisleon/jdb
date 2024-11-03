package jdb

type LinqWhere struct {
	Linq     *Linq
	A        *interface{}
	Operator string
	B        *interface{}
}

func (s *LinqWhere) And() *LinqWhere {
	return s
}
