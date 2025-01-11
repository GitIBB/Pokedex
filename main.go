package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func main() {
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
				if err := cmd.callback(); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}

	}
}

func cleanInput(text string) []string {
	lowerString := strings.ToLower(text)
	sliceString := strings.Fields(lowerString)
	return sliceString

}
