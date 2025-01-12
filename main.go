package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GitIBB/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

// Config struct for API url pagination
type Config struct {
	Next          *string
	Previous      *string
	pokeapiClient *pokeapi.Client
}

var commands map[string]cliCommand

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

// Command for showing next map page
func commandMap(cfg *Config) error {
	// Set Url
	locationResp, err := cfg.pokeapiClient.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	// Update config
	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	// Print locations
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

// Command for showing previous map page
func commandMapb(cfg *Config) error {
	if cfg.Previous == nil {
		fmt.Println("You are on the first page")
		return nil
	}

	locationResp, err := cfg.pokeapiClient.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}
	// Update config
	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	// Print locations
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil

}

func main() {
	// Create config
	cfg := &Config{}
	// Initialize command map
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
	}
	// Initialize scanner variable
	scanner := bufio.NewScanner(os.Stdin)
	// infinite for loop
	for {
		fmt.Print("Pokedex > ")
		// handle error cases for scanner
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
				if err := cmd.callback(cfg); err != nil {
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
