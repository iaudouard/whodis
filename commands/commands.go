package commands

type Command int

const (
	Exit Command = iota
	Set
	Get
	Delete
	Update
)

func (c Command) String() string {
	return [...]string{"exit", "set", "get", "delete", "update"}[c]
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
	case "update":
		return Update
	default:
		return -1
	}
}

func (c Command) IsValid() bool {
	return c != -1
}
