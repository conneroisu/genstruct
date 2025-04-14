package pkg

import (
	"time"
)

// Animal represents an animal in a zoo
type Animal struct {
	ID           string    // Unique identifier for the animal
	Name         string    // Name of the animal
	Species      string    // Species of the animal
	DateOfBirth  time.Time // Date when the animal was born
	Diet         string    // Diet type (e.g. Carnivore, Herbivore)
	Weight       float64   // Weight in kilograms
	IsEndangered bool      // Whether the species is endangered
	Habitat      string    // Natural habitat type
	Region       string    // Geographic region of origin
}

var Animals = []Animal{
	{
		ID:           "lion-001",
		Name:         "Leo",
		Species:      "African Lion",
		DateOfBirth:  time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		Diet:         "Carnivore",
		Weight:       180.5,
		IsEndangered: true,
		Habitat:      "Savanna",
		Region:       "Africa",
	},
	{
		ID:           "elephant-002",
		Name:         "Ellie",
		Species:      "African Elephant",
		DateOfBirth:  time.Date(2012, time.June, 22, 0, 0, 0, 0, time.UTC),
		Diet:         "Herbivore",
		Weight:       3200.75,
		IsEndangered: false,
		Habitat:      "Savanna",
		Region:       "Africa",
	},
	{
		ID:           "tiger-003",
		Name:         "Stripes",
		Species:      "Bengal Tiger",
		DateOfBirth:  time.Date(2019, time.February, 8, 0, 0, 0, 0, time.UTC),
		Diet:         "Carnivore",
		Weight:       160.3,
		IsEndangered: true,
		Habitat:      "Tropical Forest",
		Region:       "Asia",
	},
	{
		ID:           "penguin-004",
		Name:         "Penny",
		Species:      "Humboldt Penguin",
		DateOfBirth:  time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC),
		Diet:         "Carnivore",
		Weight:       4.2,
		IsEndangered: false,
		Habitat:      "Coastal",
		Region:       "South America",
	},
	{
		ID:           "giraffe-005",
		Name:         "George",
		Species:      "Giraffe",
		DateOfBirth:  time.Date(2016, time.April, 12, 0, 0, 0, 0, time.UTC),
		Diet:         "Herbivore",
		Weight:       1100.0,
		IsEndangered: false,
		Habitat:      "Savanna",
		Region:       "Africa",
	},
}
