package console

type printer interface {
	Print() error
}

type controller interface {
	Command() error
}
