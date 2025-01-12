package pokeapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/GitIBB/pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// parsing of API response
type LocationAreaResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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

func (c *Client) GetLocationAreas(pageURL *string) (LocationAreaResp, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if cached, ok := c.cache.Get(url); ok {
		var locResponse LocationAreaResp
		err := json.Unmarshal(cached, &locResponse)
		if err != nil {
			return LocationAreaResp{}, err
		}
		return locResponse, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationAreaResp{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	// add to cache before parsing
	c.cache.Add(url, body)

	var LocationResp LocationAreaResp
	err = json.Unmarshal(body, &LocationResp)
	if err != nil {
		return LocationAreaResp{}, err
	}

	return LocationResp, nil
}
