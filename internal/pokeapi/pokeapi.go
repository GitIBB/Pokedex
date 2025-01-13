package pokeapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/GitIBB/pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// LocationAreaResp represents the paginated list response from the PokeAPI location-areas endpoint
// count is total number of loc areas
// next and previous are urls for pagination (may be null, hence pointer type *string)
// results contains an array of location areas with names and URLs
type LocationAreaResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// Locationarea represent JSON structure returned by PokeAPI location-area endpoint
// contains area name and list of pokemon in area
// name and url for each Pokemon for additional details
type LocationArea struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Pokemon represents JSON structure returned by PokeAPI pokemon endpoint
// contains pokemon name, baseexperience and type
type Pokemon struct {
	Name string `json:"name"`

	BaseExperience int64 `json:"base_experience"`
	Height         int64 `json:"height"`
	Weight         int64 `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

type Client struct {
	// Empty for now
	cache *pokecache.Cache
}

func NewClient() *Client {
	return &Client{
		cache: pokecache.NewCache(5 * time.Second),
	}
}

// ---------------------------------------------------------

func (c *Client) GetLocationAreas(pageURL *string) (LocationAreaResp, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	// attempt to retrieve data from cache using URL as key
	// if found in cache, deserialize(unmarshal) the JSON data into LocationAreaResp struct
	// and return it, avoiding a HTTP request
	if cached, ok := c.cache.Get(url); ok {
		var locResponse LocationAreaResp
		err := json.Unmarshal(cached, &locResponse)
		if err != nil {
			return LocationAreaResp{}, err
		}
		return locResponse, nil
	}
	// Make HTTP GET request to specified URL
	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, err
	}

	defer resp.Body.Close()

	// read response body into a byte slice
	// ioutil (io) reads all bytes until end of file
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	// add to cache before parsing
	c.cache.Add(url, body)

	// declare a variable of type LocationAreaResp to store the data
	// deserialize (unmarshal) the JSON body into struct - & used as a pointer is needed to modify struct
	var LocationResp LocationAreaResp
	err = json.Unmarshal(body, &LocationResp)
	if err != nil {
		return LocationAreaResp{}, err
	}
	// return populated struct (if successful)
	return LocationResp, nil
}

// ---------------------------------------------------------

func (c *Client) GetLocationArea(name string) (LocationArea, error) {
	//url set to baseurl + location-area + name of location
	url := baseURL + "/location-area/" + name

	// attempt to retrieve data from cache using URL as key
	// if found in cache, deserialize(unmarshal) the JSON data into LocationAreaResp struct
	// and return it, avoiding a HTTP request
	if cached, ok := c.cache.Get(url); ok {
		var location LocationArea
		err := json.Unmarshal(cached, &location)
		if err != nil {
			return LocationArea{}, err
		}
		return location, nil
	}

	// Make http GET request to specified URL
	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}
	defer resp.Body.Close()

	// read response body into a byte slice
	// ioutil (io) reads all bytes until end of file
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, err
	}

	c.cache.Add(url, body)

	// declare a variable of type LocationArea to store the data
	// deserialize (unmarshal) the JSON body into struct - & used as a pointer is needed to modify struct
	var location LocationArea
	err = json.Unmarshal(body, &location)
	if err != nil {
		return LocationArea{}, err
	}

	return location, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	// attempt to retrieve data from cache using URL as key
	// if found in cache, deserialize(unmarshal) the JSON data into Pokemon struct
	// and return it, avoiding a HTTP request
	if cached, ok := c.cache.Get(url); ok {
		var pokemon Pokemon
		err := json.Unmarshal(cached, &pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	// Make http GET request to specified URL
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	// read response body into a byte slice
	// ioutil (io) reads all bytes until end of file
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, body)

	// declare a variable of type Pokemon to store the data
	// deserialize (unmarshal) the JSON body into struct - & used as a pointer is needed to modify struct
	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
