package circle_queue

type Set struct {
	list map[interface{}]bool
}

func NewSet() *Set {
	return &Set{
		list: make(map[interface{}]bool),
	}
}

func (s *Set) Add(v interface{}) bool {
	if !s.list[v] {
		s.list[v] = true
		return true
	}
	return false
}

func (s *Set) Del(v interface{}) {
	delete(s.list, v)
}

func (s *Set) Clear() {
	s.list = make(map[interface{}]bool)
}

func (s *Set) Len() int {
	return len(s.list)
}

func (s *Set) Has(v interface{}) bool {
	return s.list[v]
}

func (s *Set) Elements() []interface{} {
	ele := make([]interface{}, s.Len())
	index := 0
	for val, ok := range s.list {
		if ok {
			ele[index] = val
			index++
		}
	}
	return ele[:index]
}
