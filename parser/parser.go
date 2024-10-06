package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args string
}

func ParseCommand(input string) (*Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, fmt.Errorf("ERROR: Empty command")
	}

	return &Command{
		Name: strings.ToUpper(parts[0]),
		Args: strings.Join(parts[1:], " "),
	}, nil
}

func (c *Command) Validate() error {
	switch c.Name {
	case "SET":
		if len(strings.Fields(c.Args)) != 2 {
			return fmt.Errorf("ERROR: SET requires key and value")
		}
	case "GET", "DEL":
		if len(strings.Fields(c.Args)) != 1 {
			return fmt.Errorf("ERROR: %s requires key", c.Name)
		}
	}
	return nil
}
