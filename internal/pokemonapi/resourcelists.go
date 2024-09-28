package pokemonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const basePokemonApiUrlFormat = "https://pokeapi.co/api/v2/%v/?offset=0&limit=20"

var pageStateMap map[string]ResourceList

type ResourceList struct {
	Count    int             `json:"count"`
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	Results  []NamedResource `json:"results"`
}

type NamedResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func init() {
	pageStateMap = make(map[string]ResourceList)
}

func GetNextPage(resource string) (ResourceList, error) {
	return fetchPage(resource, func(r ResourceList) *string {
		return r.Next
	})
}

func GetPreviousPage(resource string) (ResourceList, error) {
	return fetchPage(resource, func(r ResourceList) *string {
		return r.Previous
	})
}

func fetchPage(resource string, pageRef func(ResourceList) *string) (ResourceList, error) {
	lastState, found := pageStateMap[resource]
	url := ""

	if !found {
		url = fmt.Sprintf(basePokemonApiUrlFormat, resource)
	} else {
		val := pageRef(lastState)
		if val == nil {
			return ResourceList{}, fmt.Errorf("reached the end of the pages")
		}
		url = *(pageRef(lastState))
	}
	response, err := getResourceInternal(url)
	if err != nil {
		return ResourceList{}, err
	}
	pageStateMap[resource] = response
	return response, nil
}

func getResourceInternal(url string) (ResourceList, error) {
	zero := ResourceList{}
	res, err := http.Get(url)
	fmt.Println("Calling GET " + url)
	if err != nil {
		return zero, fmt.Errorf("error when fetching %v: %w", url, err)
	}
	if res.StatusCode >= 300 {
		return zero, fmt.Errorf("error unexpected status when GET '%v': %v (%d)", url, res.Status, res.StatusCode)
	}
	if res.StatusCode == 404 {
		return zero, fmt.Errorf("resource doesn't exist at: %v", url)
	}

	resourceList := ResourceList{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resourceList)
	defer res.Body.Close()
	if err != nil {
		return zero, fmt.Errorf("error when decoding data at %v: %w", url, err)
	}

	for _, resource := range resourceList.Results {
		fmt.Printf("%v\n", resource.Name)
	}
	return resourceList, nil
}
