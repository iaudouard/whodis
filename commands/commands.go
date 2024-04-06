package commands

import "fmt"

type Command int

const (
	Exit Command = iota
	Set
	Get
	Delete
)

func (c Command) String() string {
	return [...]string{"exit", "set", "get", "delete"}[c]
}

func ParseCommand(s string) Command {
	switch s {
	case "exit":
		return Exit
	case "set":
		return Set
	case "get":
		return Get
	case "delete":
		return Delete
	default:
		return -1
	}
}

func (c Command) numberOfArgs() int {
	switch c {
	case Get:
		return 1
	case Set:
		return 2
	case Delete:
		return 1
	default:
		return 0
	}
}

func (c Command) ValidateArgs(args []string) string {
	res := ""
	if len(args) != c.numberOfArgs() {
		res = fmt.Sprintf("expected %d arguments, got %d", c.numberOfArgs(), len(args))
	}
	return res
}

func (c Command) IsValid() bool {
	return c != -1
}
