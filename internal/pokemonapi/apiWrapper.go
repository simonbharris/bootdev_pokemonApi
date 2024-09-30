package pokemonapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BaseUrl = "https://pokeapi.co/api/v2/"

// `route` should include everything after the base url, including the query parameters.
func get(route string) (*http.Response, error) {
	fullUrl := BaseUrl + route

	res, err := http.Get(fullUrl)
	fmt.Println("Calling GET " + fullUrl)
	if err != nil {
		return nil, fmt.Errorf("error when fetching %v: %w", fullUrl, err)
	}
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("resource doesn't exist at: %v", fullUrl)
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("error unexpected status when GET '%v': %v (%d)", fullUrl, res.Status, res.StatusCode)
	}
	return res, nil
}

// Retrieves the content from pokemon API at the `routeFormat` specified. It should not begin with a `/`.
// Use the %v formatter to specify where the `id` should go in the route.
// Returns a []byte of the response content.
func GetApiContentWithId(routeFormat, id string) ([]byte, error) {
	if len(routeFormat) == 0 {
		return nil, fmt.Errorf("url cannot be empty")
	}
	if len(routeFormat) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("parameters cannot be empty strings")
	}
	if !strings.Contains(routeFormat, "%v") {
		return nil, fmt.Errorf("routeFormat requires a `%%v` in the string")
	}

	fullResourceRoute := fmt.Sprintf(routeFormat, id)
	return GetApiContent(fullResourceRoute)
}

// Retrieves the content from pokemon API at the `route` specified. It should not begin with a `/`. Returns a []byte of the response content.
func GetApiContent[T struct{}](route string) ([]byte, error) {
	if len(route) == 0 {
		return nil, fmt.Errorf("url cannot be empty")
	}
	if cacheData, found := cache.Get(route); found {
		buf := &bytes.Buffer{}
		_, err := buf.Write(cacheData)
		if err != nil {
			// This shouldn't ever happen
			panic(err)
		}
		return cacheData, nil
	}

	res, err := get(route)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(route, bodyBytes)

	return bodyBytes, nil
}
