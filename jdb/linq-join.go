package jdb

type LinqJoin struct {
	Linq  *Linq
	A     *LinqFrom
	B     *LinqFrom
	Where *LinqWhere
}

/**
* On
* @param col interface{}
* @return *Linq
**/
func (s *LinqJoin) On(col interface{}) *LinqJoin {
	return s
}

func (s *LinqJoin) Eq(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) Neg(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) In(val ...interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) Like(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) More(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) Less(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) MoreEq(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) LessEs(val interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) Between(val1, val2 interface{}) *Linq {
	return s.Linq
}

func (s *LinqJoin) IsNull() *Linq {
	return s.Linq
}
