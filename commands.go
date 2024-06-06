package saga

import "log"

type Command = string

type CommandValidator struct{}

func (c CommandValidator) Validate(command Command) {
	validRoomCreatorCommands := []Command{
		"new_user:set_roles_to_rooms",
	}
	for _, validCommand := range validRoomCreatorCommands {
		if validCommand == command {
			return
		}
	}
	log.Panicf("Invalid command %s, valid commands are %v", command, validRoomCreatorCommands)
}
