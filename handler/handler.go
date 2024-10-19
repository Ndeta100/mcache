package handler

import (
	"fmt"

	"github.com/Ndeta100/mcache/parser"
	"github.com/Ndeta100/mcache/store"
)

func HandleCommand(command string, cache *store.Cache) string {
	// Step 1: Parse the command
	cmd, err := parser.ParseCommand(command)
	if err != nil {
		return err.Error()
	}

	// Step 2: Validate the command
	if err := cmd.Validate(); err != nil {
		return err.Error()
	}

	// Step 3: Execute the command based on its name
	switch cmd.Name {
	case "SET":
		if len(cmd.Args) != 2 {
			return "ERROR: SET requires key and value"
		}
		cache.Set(cmd.Args[0], cmd.Args[1])
		fmt.Printf("SET command: Key = %s, Value = %s\n", cmd.Args[0], cmd.Args[1])
		return "OK"
	case "GET":
		if len(cmd.Args) != 1 {
			return "ERROR: GET requires key"
		}
		value, found := cache.Get(cmd.Args[0])
		if found {
			fmt.Printf("GET command: Key = %s, Value = %v\n", cmd.Args[0], value)
			return fmt.Sprintf("VALUE: %v", value)
		}
		fmt.Printf("GET command: Key = %s not found\n", cmd.Args[0])
		return "NULL"
	case "DEL":
		if len(cmd.Args) != 1 {
			return "ERROR: DEL requires key"
		}
		if cache.Delete(cmd.Args[0]) {
			fmt.Printf("DEL command: Key = %s deleted\n", cmd.Args[0])
			return "OK"
		}
		fmt.Printf("DEL command: Key = %s not found\n", cmd.Args[0])
		return "ERROR: Key not found"
	default:
		return fmt.Sprintf("ERROR: Unknown command '%s'", cmd.Name)
	}
}
