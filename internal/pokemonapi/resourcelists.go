package pokemonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"pokemoncli/internal/pokecache"
	"strings"
)

const defaultRouteFormat = "%v/?offset=0&limit=20"

var pageState map[string]ResourceList
var cache pokecache.PokeCache

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
	pageState = make(map[string]ResourceList)
	cache = pokecache.NewCache(5)
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
	lastState, found := pageState[resource]
	route := ""

	if !found {
		route = fmt.Sprintf(defaultRouteFormat, resource)
	} else {
		val := pageRef(lastState)
		if val == nil {
			return ResourceList{}, fmt.Errorf("reached the end of the pages")
		}
		url := *(pageRef(lastState))
		// stripping the base url as the api wrapper always adds this.
		route = strings.Replace(url, BaseUrl, "", 1)
	}
	result := ResourceList{}
	err := getResourceInternal(route, &result)
	if err != nil {
		return ResourceList{}, err
	}
	pageState[resource] = result
	return result, nil
}

func getResourceInternal[T any](route string, out *T) error {

	content, err := GetApiContent(route)
	if err != nil {
		return err
	}

	bReader := bytes.NewReader(content)
	decoder := json.NewDecoder(bReader)
	err = decoder.Decode(out)
	if err != nil {
		return fmt.Errorf("error when decoding data at %v: %w", route, err)
	}
	return nil
}
