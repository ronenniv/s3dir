package cli

import (
	"errors"
	"fmt"
)

var validArgs = map[string]int{
	"ls": -1,
	"cd": 1,
}

type FormattedCommand struct {
	Command string
	Params  []string
}

func CheckArgs(args []string) (*FormattedCommand, error) {
	if len(args) < 1 {
		return nil, errors.New("not enough arguments")
	}

	paramCount, ok := validArgs[args[0]]
	if !ok {
		return nil, fmt.Errorf("%s not a valid argument", args[0])
	}
	if paramCount > 0 {
		if len(args[1:]) != paramCount {
			return nil, fmt.Errorf("not enough argumnets for %s", args[0])
		}
	}
	fc := FormattedCommand{}
	fc.Command = args[0]
	// fc.Params = make([]string, len(args[1:]))

	fc.Params = append(fc.Params, args[1:]...)

	return &fc, nil
}
