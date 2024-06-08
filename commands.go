package saga

import "log"

type SagaStepCommand = string

const (
	NewUserSetRolesToRoomsCommand SagaStepCommand = "new_user:set_roles_to_rooms"
)

type CommandValidator struct{}

// puede ser un interface de validación
func (c CommandValidator) Validate(command SagaStepCommand) {
	validRoomCreatorCommands := []SagaStepCommand{
		"new_user:set_roles_to_rooms",
	}
	for _, validCommand := range validRoomCreatorCommands {
		if validCommand == command {
			return
		}
	}
	log.Panicf("Invalid command %s, valid commands are %v", command, validRoomCreatorCommands)
}
