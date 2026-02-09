package validator

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/FraMan97/lestodb/internal/commands"
)

func ValidateCommand(command string) (commands.CommandInterface, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, errors.New("empty command")
	}

	elements := strings.SplitN(command, " ", 4)

	commandName := strings.ToUpper(elements[0])

	if !slices.Contains(commands.CMDS, commandName) {
		return nil, errors.New("unknown command")
	}

	switch commandName {
	case commands.CMD_SET:
		if len(elements) < commands.RULES[commands.CMD_SET] {
			return nil, errors.New("usage: SET key ttl value")
		}

		ttl, err := strconv.Atoi(elements[2])
		if err != nil {
			return nil, errors.New("invalid TTL: must be a number")
		}

		return &commands.SetCommand{
			Key:   elements[1],
			Ttl:   ttl,
			Value: elements[3],
		}, nil

	case commands.CMD_GET:
		if len(elements) < commands.RULES[commands.CMD_GET] {
			return nil, errors.New("usage: GET key")
		}
		return &commands.GetCommand{Key: elements[1]}, nil

	case commands.CMD_DEL:
		if len(elements) < commands.RULES[commands.CMD_DEL] {
			return nil, errors.New("usage: DEL key")
		}
		return &commands.DelCommand{Key: elements[1]}, nil

	case commands.CMD_BACKUP:
		if len(elements) < commands.RULES[commands.CMD_BACKUP] {
			return nil, errors.New("usage: BACKUP key|ALL")
		}
		return &commands.BackupCommand{Key: elements[1]}, nil

	case commands.CMD_RESTORE:
		if len(elements) < commands.RULES[commands.CMD_RESTORE] {
			return nil, errors.New("usage: RESTORE key|ALL ttl")
		}
		ttl, err := strconv.Atoi(elements[2])
		if err != nil {
			return nil, errors.New("invalid TTL: must be a number")
		}
		return &commands.RestoreCommand{Key: elements[1], Ttl: ttl}, nil
	}

	return nil, errors.New("internal error: command not handled")
}
