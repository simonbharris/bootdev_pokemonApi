package pokemonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const BaseUrl = "https://pokeapi.co/api/v2/"

// Retrieves the content from pokemon API at the `route` specified. It should not begin with a `/`. Returns a []byte of the response content.
func GetResource[T any](route string, out *T) error {
	content, err := getApiContent(route)
	if err != nil {
		var zero T
		if strings.HasPrefix(err.Error(), "resource doesn't exist at:") {
			slog.Info("No resources found for route: " + route)
			*out = zero
			return nil
		}
		return err
	}

	if len(content) <= 2 {
		slog.Warn(fmt.Sprintf("API Response returned with little content. Perhaps a serialization issue? \n\tRoute: %v\n\t response: %v", route, string(content[:])))
	}
	bReader := bytes.NewReader(content)
	decoder := json.NewDecoder(bReader)
	err = decoder.Decode(out)
	if err != nil {
		return fmt.Errorf("error when decoding data at %v: %w", route, err)
	}
	return nil
}

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

func getApiContentWithId(routeFormat, id string) ([]byte, error) {
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
	return getApiContent(fullResourceRoute)
}

func getApiContent[T struct{}](route string) ([]byte, error) {
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
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	cache.Add(route, bodyBytes)

	return bodyBytes, nil
}
