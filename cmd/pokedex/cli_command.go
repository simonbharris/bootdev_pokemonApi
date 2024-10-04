package main

type cliCommand struct {
	name        string
	sortOrder   int
	description string
	callback    func(args ...string) error
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			sortOrder:   1,
			description: "Returns help information",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			sortOrder:   2,
			description: "exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			sortOrder:   3,
			description: `Returns a page of 20 locations, repeated calls return the next page.`,
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			sortOrder:   4,
			description: "Returns a prior page of 20 locations. repeated calls move return prior pages.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			sortOrder:   5,
			description: "Returns what pokemon can be found in specified area.\n\tUsage: explore area-name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			sortOrder:   6,
			description: "Attempts to catch a pokemon. Some pokemon are harder to catch than others! If successful, adds them to the pokedex.",
			callback:    commandCatch,
		},
	}
}
