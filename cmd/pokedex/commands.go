package main

import (
	"fmt"
	"os"
	"pokemoncli/internal/pokemonapi"
	"sort"
)

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandHelp() error {
	cliCommands := getCliCommands()
	fmt.Println("Usage:")
	keys := make([]string, len(cliCommands))
	i := 0
	for k := range cliCommands {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, k := range keys {
		fmt.Printf("%s: %s\n", cliCommands[k].name, cliCommands[k].description)
	}
	fmt.Println()
	return nil
}

func commandMap() error {
	locations, err := pokemonapi.GetNextPage("location")
	if err != nil {
		return fmt.Errorf("error when fetching location data: %w", err)
	}

	for _, location := range locations.Results {
		fmt.Printf("%v\n", location.Name)
	}
	return nil
}

func commandMapb() error {
	locations, err := pokemonapi.GetPreviousPage("location")
	if err != nil {
		return fmt.Errorf("error when fetching location data: %w", err)
	}

	for _, location := range locations.Results {
		fmt.Printf("%v\n", location.Name)
	}
	return nil
}

/* https://pokeapi.co/api/v2/location/{id or name}/ */
