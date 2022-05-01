package ui

type ArgParseFunc func(s string) (interface{}, error)
type OptArgParseFunc func(s []string) ([]interface{}, error)

type Command struct {
	Keys         []string
	Help         string
	Args         []ArgParseFunc
	OptionalArgs OptArgParseFunc
	Action       func(c *Control, args ...interface{}) error
}
