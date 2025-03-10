package sqlite

func (s *SqlLite) GetSerie(tag string) int64
func (s *SqlLite) NextCode(tag, prefix string) string
func (s *SqlLite) SetSerie(tag string, val int) int64
func (s *SqlLite) CurrentSerie(tag string) int64
