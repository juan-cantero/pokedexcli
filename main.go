package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/juan-cantero/pokedexcli/internal/pokeapi"
	"github.com/juan-cantero/pokedexcli/internal/pokedex"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
	config      *config
}

type config struct {
	Next     string
	Previous string
	Pokedex  *pokedex.Pokedex
}

var cmds map[string]*cliCommand

func main() {
	mapConfig := &config{
		Pokedex: pokedex.New(),
	}
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
		"explore": {
			name:        "explore",
			description: "explore area",
			callback:    commandExplore,
			config:      mapConfig,
		},
		"catch": {
			name:        "catch",
			description: "catch pokemon",
			callback:    commandCatch,
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
			// Pass arguments (everything after the command name)
			args := words[1:]
			err := cmd.callback(cmd.config, args)
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

func commandMapForward(cfg *config, args []string) error {
	url := cfg.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	}
	return fetchAndUpdateConfig(cfg, url)
}

func commandMapBackward(cfg *config, args []string) error {
	if cfg.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	}
	return fetchAndUpdateConfig(cfg, cfg.Previous)
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore requires an area name. Usage: explore <area-name>")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	area, err := pokeapi.FetchPokemonsByArea(areaName)
	if err != nil {
		return fmt.Errorf("failed to explore area: %w", err)
	}
	fmt.Printf("Exploring %s...", areaName)
	fmt.Println()
	fmt.Println("Found Pokemons:")
	names := area.GetPokemonNames()
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch requires a pokemon name. Usage: catch <pokemon-name>")
	}

	pokemonName := args[0]

	// Check if already caught
	if cfg.Pokedex.Has(pokemonName) {
		fmt.Printf("You already have %s in your Pokedex!\n", pokemonName)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := pokeapi.FetchPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	// Calculate catch probability based on base experience
	// Higher base experience = harder to catch
	// Base experience ranges from ~50-600+
	// We'll make it so lower base exp = higher success rate
	// Success rate: 100 - (base_exp / 3), capped between 20% and 80%
	successRate := max(20, min(80, 100-(pokemon.BaseExperience/3)))

	// Roll a random number between 0-100
	roll := rand.Intn(100)

	// If roll is greater than success rate, catch fails
	if roll > successRate {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	// Successfully caught!
	cfg.Pokedex.Catch(*pokemon)
	fmt.Printf("%s was caught!\n", pokemon.Name)

	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, v := range cmds {
		fmt.Printf("%v: %v", v.name, v.description)
		fmt.Println()
	}
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
