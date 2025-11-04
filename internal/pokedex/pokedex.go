package pokedex

import (
	"fmt"

	"github.com/juan-cantero/pokedexcli/internal/models"
)

// Pokedex manages the collection of caught Pokemon
type Pokedex struct {
	pokemon map[string]models.Pokemon
}

// New creates a new Pokedex instance
func New() *Pokedex {
	return &Pokedex{
		pokemon: make(map[string]models.Pokemon),
	}
}

// Catch adds a Pokemon to the Pokedex
// Returns true if the Pokemon was newly caught, false if already in the Pokedex
func (p *Pokedex) Catch(pokemon models.Pokemon) bool {
	if _, exists := p.pokemon[pokemon.Name]; exists {
		return false
	}
	p.pokemon[pokemon.Name] = pokemon
	return true
}

// Has checks if a Pokemon is already in the Pokedex
func (p *Pokedex) Has(name string) bool {
	_, exists := p.pokemon[name]
	return exists
}

// Get retrieves a Pokemon from the Pokedex
func (p *Pokedex) Get(name string) (models.Pokemon, error) {
	pokemon, exists := p.pokemon[name]
	if !exists {
		return models.Pokemon{}, fmt.Errorf("pokemon %s not found in pokedex", name)
	}
	return pokemon, nil
}

// List returns all caught Pokemon
func (p *Pokedex) List() []models.Pokemon {
	list := make([]models.Pokemon, 0, len(p.pokemon))
	for _, pokemon := range p.pokemon {
		list = append(list, pokemon)
	}
	return list
}

// Count returns the number of caught Pokemon
func (p *Pokedex) Count() int {
	return len(p.pokemon)
}
