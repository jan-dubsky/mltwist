package expr

type Ref struct {
	s string
}

func NewRef(s string) Ref {
	return Ref{s: s}
}
