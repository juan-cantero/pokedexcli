package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/juan-cantero/pokedexcli/internal/models"
	"github.com/juan-cantero/pokedexcli/internal/pokecache"
)

// Global cache instance
var cache *pokecache.Cache

func init() {
	// Initialize cache with 5 minute cleanup interval
	cache = pokecache.NewCache(5 * time.Minute)
}

func FetchLocationAreas(url string) (*models.LocationArea, error) {
	// Check if data is in cache
	if cachedData, found := cache.Get(url); found {
		// Cache hit! Unmarshal and return
		var la models.LocationArea
		if err := json.Unmarshal(cachedData, &la); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return &la, nil
	}

	// Cache miss - fetch from API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Store in cache for next time
	cache.Add(url, body)

	// Unmarshal and return
	var la models.LocationArea
	if err := json.Unmarshal(body, &la); err != nil {
		return nil, err
	}

	return &la, nil
}

func FetchPokemonsByArea(areaName string) (*models.PokemonsByArea, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	if cachedData, found := cache.Get(url); found {
		var pba models.PokemonsByArea
		if err := json.Unmarshal(cachedData, &pba); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data")
		}
		return &pba, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status")
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body")
	}
	cache.Add(url, body)

	var pba models.PokemonsByArea
	if err := json.Unmarshal(body, &pba); err != nil {
		return nil, err
	}
	return &pba, nil

}

func FetchPokemon(pokemonName string) (*models.Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	// Check cache first
	if cachedData, found := cache.Get(url); found {
		var pokemon models.Pokemon
		if err := json.Unmarshal(cachedData, &pokemon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return &pokemon, nil
	}

	// Fetch from API
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Cache the response
	cache.Add(url, body)

	// Unmarshal and return
	var pokemon models.Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pokemon data: %w", err)
	}

	return &pokemon, nil
}
