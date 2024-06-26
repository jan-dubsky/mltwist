package consoleui

type ArgParseFunc func(s string) (interface{}, error)
type OptArgParseFunc func(s []string) ([]interface{}, error)

type Command struct {
	Keys         []string
	Help         string
	Args         []ArgParseFunc
	OptionalArgs OptArgParseFunc
	Action       func(ui *UI, args ...interface{}) error
}
