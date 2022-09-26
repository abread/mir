package aba

type void struct{}
type Set[T comparable] map[T]void

func (s Set[T]) Add(value T) {
	s[value] = void{}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Has(value T) bool {
	_, present := s[value]
	return present
}

func (s Set[T]) Len() int {
	return len(s)
}
