package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"gopkg.in/yaml.v3"
)

// ItemType represents a type of item that can be placed in the hex grid
type ItemType struct {
	Name       string  `yaml:"name"`
	Percentage float64 `yaml:"percentage"`
	Style      string  `yaml:"style"` // "dot" or "fill"
	Color      string  `yaml:"color"`
}

// Config represents the YAML configuration file structure
type YAMLConfig struct {
	Default string     `yaml:"default"`
	Items   []ItemType `yaml:"items"`
}

// HexCell represents a single hexagon cell in the grid
type HexCell struct {
	Row      int
	Col      int
	ItemType *ItemType
	X, Y     float64 // Center coordinates
}

// HexGrid represents the complete hex grid
type HexGrid struct {
	Rows         int
	Cols         int
	Cells        [][]*HexCell
	ItemTypes    []*ItemType
	DefaultColor string
}

// LoadYAMLConfig loads and parses the YAML configuration file
func LoadYAMLConfig(filePath string) (*YAMLConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	var config YAMLConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate configuration
	if config.Default == "" {
		return nil, fmt.Errorf("default color is required")
	}
	if len(config.Items) == 0 {
		return nil, fmt.Errorf("no items defined in configuration")
	}

	totalPercentage := 0.0
	for _, item := range config.Items {
		if item.Percentage < 0 || item.Percentage > 100 {
			return nil, fmt.Errorf("invalid percentage for item %s: %f", item.Name, item.Percentage)
		}
		totalPercentage += item.Percentage
		if item.Style != "dot" && item.Style != "fill" {
			return nil, fmt.Errorf("invalid style for item %s: %s (must be 'dot' or 'fill')", item.Name, item.Style)
		}
	}

	if totalPercentage > 100 {
		return nil, fmt.Errorf("total percentage exceeds 100%%: %f", totalPercentage)
	}

	return &config, nil
}

// CreateHexGrid creates a new hex grid with the specified dimensions
func CreateHexGrid(rows, cols int, config *YAMLConfig) *HexGrid {
	grid := &HexGrid{
		Rows:         rows,
		Cols:         cols,
		Cells:        make([][]*HexCell, rows),
		ItemTypes:    make([]*ItemType, len(config.Items)),
		DefaultColor: config.Default,
	}

	// Copy item types
	for i := range config.Items {
		grid.ItemTypes[i] = &config.Items[i]
	}

	// Initialize cells
	for row := 0; row < rows; row++ {
		grid.Cells[row] = make([]*HexCell, cols)
		for col := 0; col < cols; col++ {
			grid.Cells[row][col] = &HexCell{
				Row: row,
				Col: col,
			}
		}
	}

	return grid
}

// PopulateGrid fills the grid with items based on their percentages
func (grid *HexGrid) PopulateGrid() {
	totalCells := grid.Rows * grid.Cols

	// Calculate how many cells each item type should occupy
	itemCounts := make(map[*ItemType]int)
	for _, itemType := range grid.ItemTypes {
		count := int(float64(totalCells) * itemType.Percentage / 100.0)
		itemCounts[itemType] = count
	}

	// Create a list of all cells
	allCells := make([]*HexCell, 0, totalCells)
	for row := 0; row < grid.Rows; row++ {
		for col := 0; col < grid.Cols; col++ {
			allCells = append(allCells, grid.Cells[row][col])
		}
	}

	// Shuffle the cells
	rand.Shuffle(len(allCells), func(i, j int) {
		allCells[i], allCells[j] = allCells[j], allCells[i]
	})

	// Assign items to cells
	cellIndex := 0
	for itemType, count := range itemCounts {
		for i := 0; i < count && cellIndex < len(allCells); i++ {
			allCells[cellIndex].ItemType = itemType
			cellIndex++
		}
	}
}
