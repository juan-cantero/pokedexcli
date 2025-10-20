package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/juan-cantero/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
	config      *config
}

type config struct {
	Next     string
	Previous string
}

var cmds map[string]*cliCommand

func main() {
	mapConfig := &config{}
	cmds = map[string]*cliCommand{
		"map": {
			name:        "map",
			description: "show the next page of locations",
			callback:    commandMapForward,
			config:      mapConfig,
		},
		"mapb": {
			name:        "mapb",
			description: "show the previous page of locations",
			callback:    commandMapBackward,
			config:      mapConfig,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleaned := strings.ToLower(strings.TrimSpace(input))
		words := strings.Fields(cleaned)

		// Handle empty input
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		// Look up command in registry
		if cmd, exists := cmds[commandName]; exists {
			err := cmd.callback(cmd.config)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	return words
}

func fetchAndUpdateConfig(cfg *config, url string) error {
	location, err := pokeapi.FetchLocationAreas(url)
	if err != nil {
		return err
	}

	// Print results
	for _, l := range location.Results {
		fmt.Println(l.Name)
	}

	// Update config
	cfg.Next = location.Next
	if location.Previous == nil {
		cfg.Previous = ""
	} else if s, ok := location.Previous.(string); ok {
		cfg.Previous = s
	} else {
		cfg.Previous = fmt.Sprintf("%v", location.Previous)
	}

	return nil
}

func commandMapForward(cfg *config) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}
	return fetchAndUpdateConfig(cfg, url)
}

func commandMapBackward(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	return fetchAndUpdateConfig(cfg, cfg.Previous)
}

func commandHelp(cfg *config) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, v := range cmds {
		fmt.Printf("%v: %v", v.name, v.description)
		fmt.Println()
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
