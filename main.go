package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/GitIBB/pokedex/internal/pokeapi"
)

// cliCommand represents a single command in the Pokedex CLI
type cliCommand struct {
	name        string                         // cmd identifier
	description string                         // help text for command
	callback    func(*Config, ...string) error // func to exectue when command is called
}

// holds state and dependencies for pokedex application
type Config struct {
	Next          *string         //url for next page
	Previous      *string         // url for previous page
	pokeapiClient *pokeapi.Client // client for making PokeAPI requests
	caughtPokemon map[string]pokeapi.Pokemon
}

// storage of all available CLI commands mapped by names
var commands map[string]cliCommand

// handles exit command, cleaning up and terminating the program
func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// displays all available commands and their descriptions
func commandHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

// Command which retrieves and displays next page of location areas from PokeAPI, updating pagination state in config
func commandMap(cfg *Config, args ...string) error {
	// Fetch next page of location areas using the stored url
	locationResp, err := cfg.pokeapiClient.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	// Update pagination URLs for next/previous nav
	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	// Display all location areas from current page
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

// Same as command, but for previous page of location areas
func commandMapb(cfg *Config, args ...string) error {
	// check if previous page exists
	if cfg.Previous == nil {
		fmt.Println("You are on the first page")
		return nil
	}
	// fetch previous page of loc areas
	locationResp, err := cfg.pokeapiClient.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}
	// Update pagination URLs for nav
	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	// Print locations from current page
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil

}

// commandExplore retrieves and displays all Pokemon that can be encountered
// in a specific location area - takes one argument (location area's name)
func commandExplore(cfg *Config, args ...string) error {
	// validate that exactly one location area (arg) name was provided
	if len(args) != 1 {
		return errors.New("location area name is required")
	}
	locationAreaName := args[0]
	// Fetch detailed information about specific location area
	location, err := cfg.pokeapiClient.GetLocationArea(locationAreaName)
	if err != nil {
		return err
	}

	// Display all Pokemon that can be encountered in this area
	fmt.Printf("Exploring %s...\n", locationAreaName)
	fmt.Println("Found Pokemnon:")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("pokemon name is required")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}
	if _, ok := cfg.caughtPokemon[pokemonName]; ok {
		return errors.New("Already caught this pokemon!")
	}

	// formula for probability of catching pokemon
	baseRate := 100
	catchRate := baseRate - int(pokemon.BaseExperience/4)
	if catchRate < 5 {
		catchRate = 5
	}

	randomNum := rand.Intn(baseRate)

	if randomNum < catchRate {
		cfg.caughtPokemon[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)

	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)

	}
	return nil
}

func commandInspect(cfg *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("pokemon name is required")
	}

	name := args[0]

	if pokemon, ok := cfg.caughtPokemon[name]; !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}

	}
	return nil

}

func commandPokedex(cfg *Config, args ...string) error {

	for key, _ := range cfg.caughtPokemon {
		fmt.Println(key)
	}
	return nil
}

// Initializes and runs Pokedex CLI application
// sets up config, available commands and starts the REPL
func main() {
	// initialize config with new PokeAPI client
	cfg := &Config{
		pokeapiClient: pokeapi.NewClient(),
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}
	// Initialize command map with command implementation
	commands = map[string]cliCommand{
		"exit": {
			// Exit command
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		// Help command
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		// Map command
		"map": {
			name:        "map",
			description: "Shows the next 20 location areas",
			callback:    commandMap,
		},
		// Map(back) command
		"mapb": {
			name:        "mapb",
			description: "Shows the previous 20 location areas",
			callback:    commandMapb,
		},
		// explore command, lists the pokemon present in location area
		"explore": {
			name:        "explore",
			description: "Explores area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View Pokedex",
			callback:    commandPokedex,
		},
	}
	// Initialize scanner variable (Start the REPL, r)
	scanner := bufio.NewScanner(os.Stdin)
	// infinite for loop
	for {
		fmt.Print("Pokedex > ")
		// Read user input, break loop on EOF or error
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
			break
		}
		// Process the scanned text and clean it
		text := scanner.Text()
		words := cleanInput(text)

		if len(words) > 0 {
			commandKey := words[0] // Assuming the first word is the command

			// Use commandKey to look up the command
			if cmd, exists := commands[commandKey]; exists {
				args := words[1:]
				if err := cmd.callback(cfg, args...); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}

	}
}

// Cleans terminal input
func cleanInput(text string) []string {
	lowerString := strings.ToLower(text)
	sliceString := strings.Fields(lowerString)
	return sliceString

}
