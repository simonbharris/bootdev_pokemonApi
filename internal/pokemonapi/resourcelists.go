package pokemonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokemoncli/internal/pokecache"
)

const basePokemonApiUrlFormat = "https://pokeapi.co/api/v2/%v/?offset=0&limit=20"

var pageStateMap map[string]ResourceList
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
	pageStateMap = make(map[string]ResourceList)
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

func getUrl(url string) (*http.Response, error) {
	res, err := http.Get(url)
	fmt.Println("Calling GET " + url)
	if err != nil {
		return nil, fmt.Errorf("error when fetching %v: %w", url, err)
	}
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("resource doesn't exist at: %v", url)
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("error unexpected status when GET '%v': %v (%d)", url, res.Status, res.StatusCode)
	}
	return res, nil
}

func getResourceInternal(url string) (ResourceList, error) {
	zero := ResourceList{}
	resourceList := ResourceList{}

	if cacheData, found := cache.Get(url); found {
		buf := &bytes.Buffer{}
		_, err := buf.Write(cacheData)
		if err != nil {
			// This shouldn't ever happen
			panic(err)
		}

		decoder := json.NewDecoder(buf)
		err = decoder.Decode(&resourceList)
		fmt.Printf("%v Was found in cache\n", url)
		if err != nil {
			return ResourceList{}, err
		}
		return resourceList, nil
	}

	res, err := getUrl(url)
	if err != nil {
		return ResourceList{}, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return ResourceList{}, err
	}
	cache.Add(url, bodyBytes)

	bReader := bytes.NewReader(bodyBytes)
	decoder := json.NewDecoder(bReader)
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
