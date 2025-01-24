package main

import (
	"bufio"
	"fmt"
	"github.com/kjittiwan/pokedexcli/internal/pokeapi"
	"github.com/kjittiwan/pokedexcli/internal/pokecache"
	"math/rand"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *pokeapi.LocationAreasConfig, cache *pokecache.Cache, parameter string, pokedex Pokedex) error
}

type Pokedex map[string]pokeapi.Pokemon

var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area, usage: explore <area_name>",
			callback:    commandExplore},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon, usage: catch <pokemon>",
			callback:    commandCatch},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon in your pokedex",
			callback:    commandInspect},
		"pokedex": {
			name:        "pokedex",
			description: "See pokemons you've caught in your pokedex",
			callback:    commandPokedex},
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	locationsAreaConfig := pokeapi.LocationAreasConfig{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: nil,
	}
	cache := pokecache.NewCache(1 * time.Minute)
	pokedex := map[string]pokeapi.Pokemon{}
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		inputSlice := cleanInput(scanner.Text())
		command, ok := supportedCommands[inputSlice[0]]
		var parameter string
		if len(inputSlice) > 1 {
			parameter = inputSlice[1]

		}
		if !ok {
			fmt.Println("Unknown command")
		} else {
			err := command.callback(&locationsAreaConfig, cache, parameter, pokedex)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func commandExit(config *pokeapi.LocationAreasConfig, cache *pokecache.Cache, _ string, _ Pokedex) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.LocationAreasConfig, cache *pokecache.Cache, _ string, _ Pokedex) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println(" ")
	for _, val := range supportedCommands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandMap(config *pokeapi.LocationAreasConfig, cache *pokecache.Cache, _ string, _ Pokedex) error {
	locationAreas, err := pokeapi.GetLocationsArea(config.Next, cache)
	if err != nil {
		return fmt.Errorf("error executing map command: %w", err)
	}

	// loop over all areas and print
	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	// set next and prev
	config.Next = locationAreas.Next
	config.Previous = locationAreas.Previous
	return nil
}
func commandMapB(config *pokeapi.LocationAreasConfig, cache *pokecache.Cache, _ string, _ Pokedex) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := pokeapi.GetLocationsArea(*config.Previous, cache)
	if err != nil {
		return fmt.Errorf("error executing mapb command: %w", err)
	}

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	config.Next = locationAreas.Next
	config.Previous = locationAreas.Previous
	return nil
}

func commandExplore(_ *pokeapi.LocationAreasConfig, cache *pokecache.Cache, area string, _ Pokedex) error {
	if area == "" {
		fmt.Println("Please specify an area to explore")
		return nil
	}
	fmt.Printf("Exploring %s...\n", area)
	locationInfo, err := pokeapi.GetLocationInfo(area, cache)
	if err != nil {
		return fmt.Errorf("error executing explore command: %w", err)
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationInfo.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(_ *pokeapi.LocationAreasConfig, _ *pokecache.Cache, pokemonName string, pokedex Pokedex) error {
	if pokemonName == "" {
		fmt.Println("Please specify a Pokemon to catch")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("error executing command catch: %w", err)
	}
	// delay for suspense
	time.Sleep(time.Second * 1)
	// base_ex < 50 -> 90%
	// base_ex < 150 -> 75%
	// base_ex < 250 -> 70%
	// base_ex > 250 -> 60%
	var probability float64
	baseExp := pokemon.BaseExperience
	switch {
	case baseExp < 50:
		probability = 0.9
	case baseExp < 150:
		probability = 0.75
	case baseExp < 250:
		probability = 0.7
	case baseExp > 300:
		probability = 0.6
	default:
		probability = 0.8
	}
	if rand.Float64() < probability {
		pokedex[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func commandInspect(_ *pokeapi.LocationAreasConfig, _ *pokecache.Cache, pokemonName string, pokedex Pokedex) error {
	// check if exists in pokedex,
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		fmt.Println("You haven't caught this pokemon yet!")
		return nil
	}
	pokeapi.PrintPokemon(pokemon)
	return nil
}

func commandPokedex(_ *pokeapi.LocationAreasConfig, _ *pokecache.Cache, pokemonName string, pokedex Pokedex) error {
	if len(pokedex) == 0 {
		fmt.Println("Your pokedex is empty, better catch some pokemons!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for key := range pokedex {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
