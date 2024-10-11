package main

import (
	"fmt"
	"os"
	"pokemoncli/internal/pokedexservice"
	"pokemoncli/internal/pokemonapi"
	"sort"
	"strconv"
)

var pokedex map[string]pokemonapi.Pokemon

func init() {
	pokedex = make(map[string]pokemonapi.Pokemon)
}

func commandExit(args ...string) error {
	os.Exit(0)
	return nil
}

func commandHelp(args ...string) error {
	cliCommands := getCliCommands()
	fmt.Println("Usage:")
	arr := make([]cliCommand, len(cliCommands))
	i := 0
	for _, v := range cliCommands {
		arr[i] = v
		i++
	}
	sort.Slice(arr, func(i, j int) bool {
		// ascending order
		return arr[i].sortOrder < arr[j].sortOrder
	})
	for i := range arr {
		fmt.Printf("%s: %s\n", arr[i].name, arr[i].description)
	}
	fmt.Println()
	return nil
}

func commandMap(args ...string) error {
	locations, err := pokemonapi.GetNextResourceListPage("location")
	if err != nil {
		return fmt.Errorf("error when fetching location data: %w", err)
	}

	for _, location := range locations.Results {
		fmt.Printf("%v\n", location.Name)
	}
	return nil
}

func commandMapb(args ...string) error {
	locations, err := pokemonapi.GetPreviousResourceListPage("location")
	if err != nil {
		return fmt.Errorf("error when fetching location data: %w", err)
	}

	for _, location := range locations.Results {
		fmt.Printf("%v\n", location.Name)
	}
	return nil
}

func commandExplore(args ...string) error {
	if len(args) != 1 {
		fmt.Println("explore requires exactly 1 area name. See help")
		return nil
	}
	locationName := args[0]
	pokemonNames, err := pokedexservice.Explore(locationName)
	if err != nil {
		return err
	}
	if len(pokemonNames) == 0 {
		fmt.Println("Area has no pokemon.")
	}
	for _, name := range pokemonNames {
		fmt.Printf(" - %v\n", name)
	}
	return nil
}

func commandCatch(args ...string) error {
	if len(args) != 1 {
		fmt.Println("explore requires exactly 1 area name. See help")
		return nil
	}
	pokemonName := args[0]

	isCaught, err := pokedexservice.Catch(pokedex, pokemonName)
	if err != nil {
		return err
	}

	if isCaught {
		fmt.Printf("%v was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the `inspect` command.")
	} else {
		fmt.Printf("%v escaped!\n", pokemonName)
	}
	return nil
}

func commandInspect(args ...string) error {
	if len(args) != 1 {
		fmt.Println("explore requires exactly 1 area name. See help")
		return nil
	}
	pokemonName := args[0]

	pokemon := pokedexservice.Inspect(pokedex, pokemonName)

	printMainValue("Name", pokemon.Name)
	printMainValue("Height", strconv.Itoa(pokemon.Height))
	printMainValue("Weight", strconv.Itoa(pokemon.Weight))
	printMainValue("Stats", "")
	for _, statData := range pokemon.Stats {
		printListValues(statData.Stat.Name, strconv.Itoa(statData.BaseStat))
	}
	printMainValue("Types", "")
	for _, pType := range pokemon.Types {
		printListValues(pType.Type.Name, "")
	}
	return nil
}

func commandPokedex(args ...string) error {
	if len(args) > 0 {
		fmt.Println("Pokedex does not take any arguments")
		return nil
	}
	pokemonList := pokedexservice.Pokedex(pokedex)
	if len(pokedex) == 0 {
		fmt.Println("No pokemon registered in pokedex yet, go `catch` some pokemon!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, mon := range pokemonList {
		printListValues(mon, "")
	}
	return nil
}

/* https://pokeapi.co/api/v2/location/{id or name}/ */
