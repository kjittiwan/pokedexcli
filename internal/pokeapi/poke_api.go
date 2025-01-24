package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kjittiwan/pokedexcli/internal/pokecache"
)

type LocationAreas struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type LocationAreasConfig struct {
	Next     string
	Previous *string
}

type LocationInfo struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func GetLocationsArea(url string, cache *pokecache.Cache) (LocationAreas, error) {
	data, ok := cache.Get(url)
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return LocationAreas{}, err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return LocationAreas{}, err
		}
		cache.Add(url, data)
	}

	var locationAreas LocationAreas
	if err := json.Unmarshal(data, &locationAreas); err != nil {
		return LocationAreas{}, err
	}
	return locationAreas, nil
}

func GetLocationInfo(area string, cache *pokecache.Cache) (LocationInfo, error) {
	data, ok := cache.Get(area)
	if !ok {
		resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area))
		if err != nil {
			return LocationInfo{}, err
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return LocationInfo{}, err
		}
		cache.Add(area, data)
	}
	var locationInfo LocationInfo
	if err := json.Unmarshal(data, &locationInfo); err != nil {
		return LocationInfo{}, err
	}
	return locationInfo, nil
}

func GetPokemon(name string) (Pokemon, error) {
	resp, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name))
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	var pokemon Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func PrintPokemon(pokemon Pokemon) {
	fmt.Printf("Name: %d\n", pokemon.Height)
	fmt.Printf("Height: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokeType.Type.Name)
	}
}
