package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args string
}
type File struct {
	Commands []Command
}

func (f File) String() string {
	var sb strings.Builder
	for _, cmd := range f.Commands {
		fmt.Fprintln(&sb, cmd.String())
	}

	return sb.String()
}
func (c Command) String() string {
	var sb strings.Builder
	switch c.Name {
	case "model":
		fmt.Fprintf(&sb, "FROM %s", c.Args)
	case "license", "template", "system", "adapter":
		fmt.Fprintf(&sb, "%s %s", strings.ToUpper(c.Name), quote(c.Args))
	case "message":
		role, message, _ := strings.Cut(c.Args, ": ")
		fmt.Fprintf(&sb, "MESSAGE %s %s", role, quote(message))
	default:
		fmt.Fprintf(&sb, "PARAMETER %s %s", c.Name, quote(c.Args))
	}

	return sb.String()
}

func quote(s string) string {
	if strings.Contains(s, "\n") || strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
		if strings.Contains(s, "\"") {
			return `"""` + s + `"""`
		}

		return `"` + s + `"`
	}

	return s
}
