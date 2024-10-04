package pokemonapi

import (
	"fmt"
	"pokemoncli/internal/shcache"
	"strings"
)

const defaultRouteFormat = "%v/?offset=0&limit=20"

var resourceListPageState map[string]ResourceList
var cache shcache.Cache

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

func GetNextResourceListPage(resource string) (ResourceList, error) {
	return getResourceListByPage(resource, func(r ResourceList) *string {
		return r.Next
	})
}

func GetPreviousResourceListPage(resource string) (ResourceList, error) {
	return getResourceListByPage(resource, func(r ResourceList) *string {
		return r.Previous
	})
}

// Retrieves the content from pokemon API at the `routeFormat` specified. It should not begin with a `/`.
// Use the %v formatter to specify where the `id` should go in the route.
// Returns a []byte of the response content.
func GetResourceWithId[T any](routeFormat, id string, out *T) error {
	return GetResource(fmt.Sprintf(routeFormat, id), out)
}

func init() {
	resourceListPageState = make(map[string]ResourceList)
	cache = shcache.NewCache(5)
}

func getResourceListByPage(resource string, pageRef func(ResourceList) *string) (ResourceList, error) {
	lastState, found := resourceListPageState[resource]
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
	err := GetResource(route, &result)
	if err != nil {
		return ResourceList{}, err
	}
	resourceListPageState[resource] = result
	return result, nil
}
