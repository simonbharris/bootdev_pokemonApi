package pokemonapi

import (
	"fmt"
	"log/slog"
)

type LocationArea struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func Explore(locationName string) ([]string, error) {
	if len(locationName) == 0 {
		return nil, fmt.Errorf("location cannot be an empty string")
	}
	content := LocationArea{}
	err := GetResourceWithId("location-area/%v", locationName, &content)
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
