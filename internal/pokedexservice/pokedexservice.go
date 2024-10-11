package pokedexservice

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"pokemoncli/internal/pokemonapi"
	"strconv"
	"strings"
)

var discoveredPokemon map[string]pokemonapi.Pokemon

func init() {
	discoveredPokemon = make(map[string]pokemonapi.Pokemon)
}

func Inspect(pokemonName string) {
	pokemon, found := discoveredPokemon[pokemonName]
	if !found {
		fmt.Println("you have not caught that pokemon")
		return
	}
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
}

func printMainValue(key, value string) {
	fmt.Printf("%v: %v\n", key, value)
}

// if value is 0-length, don't print a `:`
func printListValues(key, value string) {
	if len(value) == 0 {
		fmt.Printf("  - %v\n", strings.ToLower(key))
	} else {
		fmt.Printf("  -%v: %v\n", strings.ToLower(key), value)
	}
}

// Attempts to capture a pokemon. If successful it will be stored internally and return true
// else returns false.
func Catch(pokemonName string) (bool, error) {
	if len(pokemonName) == 0 {
		return false, fmt.Errorf("pokemonName cannot be an empty string")
	}
	pokemon := pokemonapi.Pokemon{}
	err := pokemonapi.GetResourceWithId("pokemon/%v", pokemonName, &pokemon)
	if err != nil {
		return false, err
	}

	isCaught := calculateCapture(pokemon.BaseExperience)
	if !isCaught {
		return false, nil
	}
	if _, found := discoveredPokemon[pokemonName]; !found {
		discoveredPokemon[pokemonName] = pokemon
	}
	return true, nil
}

// Higher the difficulty, the harder it is to capture.
func calculateCapture(difficulty int) bool {
	// numerator / difficulty being the lowest success rate
	numerator := math.Pow(float64(difficulty), 0.5)
	difference := float64(difficulty) - numerator
	chance := (numerator + (rand.Float64() * difference / 5)) / float64(difficulty)
	luck := rand.Float64()
	slog.Debug(fmt.Sprintf("Difficulty: %v, numerator %v, catch chance: %v, Luck: %v, caught: %v", difficulty, numerator, chance, luck, luck > 1-chance))
	return luck > 1-chance
}

// Prints an array of capturable pokemon from a location.
// If no pokemon exist in location, returns an empty array.
func Explore(locationName string) ([]string, error) {
	if len(locationName) == 0 {
		return nil, fmt.Errorf("location cannot be an empty string")
	}
	content := pokemonapi.LocationArea{}
	err := pokemonapi.GetResourceWithId("location-area/%v", locationName, &content)
	if err != nil {
		return nil, fmt.Errorf("error when fetching location-area: %w", err)
	}
	if len(content.PokemonEncounters) == 0 {
		slog.Debug("Location-area returned with no encounters")
		return []string{}, nil
	}
	results := make([]string, len(content.PokemonEncounters))
	for i, pokemonEncounter := range content.PokemonEncounters {
		results[i] = pokemonEncounter.Pokemon.Name
	}
	return results, nil
}
