package main

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Returns help information",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: `Returns a page of 20 locations, repeated calls return the next page.`,
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: `Returns a prior page of 20 locations. repeated calls move return prior pages.`,
			callback:    commandMapb,
		},
	}
}
