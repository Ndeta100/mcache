package parser

import (
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(input string) (*Command, error) {
	// Split the input string by whitespace
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return nil, fmt.Errorf("ERROR: Empty command")
	}

	// Create and return a Command instance
	return &Command{
		Name: strings.ToUpper(parts[0]),
		Args: parts[1:], // Store arguments as a slice of strings
	}, nil
}

func (c *Command) Validate() error {
	switch c.Name {
	case "SET":
		if len(c.Args) != 2 {
			return fmt.Errorf("ERROR: SET requires key and value")
		}
	case "GET", "DEL":
		if len(c.Args) != 1 {
			return fmt.Errorf("ERROR: %s requires a key", c.Name)
		}
	default:
		return fmt.Errorf("ERROR: Unknown command '%s'", c.Name)
	}
	return nil
}
