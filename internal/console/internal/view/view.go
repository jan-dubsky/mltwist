package view

type View interface {
	// MinLines returns minimal number of lines the view needs to render.
	// The returned number should be always non-negative. Negative number
	// will be understood as zero.
	//
	// For the same state of the structure, consecutive calls of this method
	// must return the same result. Consequently, it's recommended for this
	// method not to change any internal state of the structure.
	MinLines() int
	// MaxLines returns maximal number of lines the view is able to render
	// to. The value returned should be greater or equal to value returned
	// by MinLines(). Positive values less then MinLines will be
	// interpretted as MinLines(). Negative values are be interpreted as
	// infinity.
	//
	// For the same state of the structure, consecutive calls of this method
	// must return the same result. Consequently, it's recommended for this
	// method not to change any internal state of the structure.
	MaxLines() int

	// Print instructs the element to prints EXACTLY n lines to the screen.
	Print(n int) error
}
