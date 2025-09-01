# Hex Grid Generator

A Golang application that generates hex grids from YAML configuration files with a graphical user interface.

## Features

- GUI for selecting YAML configuration files and output locations
- Configurable grid size (rows and columns)
- Support for different item types with percentages, styles, and colors
- SVG output with proper staggered hexagon grid layout (no overlapping)
- HTML page with embedded SVG, scrolling, and item legend
- Two item styles: "fill" (colored hexagon) and "dot" (colored dot in center with black outline)
- Configurable default background color for empty cells and dot-style items

## Installation

1. Make sure you have Go 1.21 or later installed
2. Clone or download this repository
3. Run `go mod tidy` to download dependencies
4. Run `go run .` to start the application

## Usage

1. **Start the application**: Run `go run .` from the project directory
2. **Select YAML file**: Choose from the dropdown menu to select your configuration file from the `grid-specs/` folder
3. **Set grid size**: Enter the number of rows and columns for your hex grid
4. **Auto-generated output**: The output path is automatically generated based on the YAML filename and timestamp
5. **Generate**: Click "Generate Hex Grid" to create the SVG and HTML files in the `generated-grids/` folder
6. **Auto-open**: The generated HTML file automatically opens in your default browser

### File Structure

The application uses the following directory structure:
- `grid-specs/` - Place your YAML configuration files here
- `generated-grids/` - Generated SVG and HTML files are saved here with automatic naming

### YAML File Selection

The application automatically scans the `grid-specs/` folder for YAML files and displays them in a dropdown menu. Use the "Refresh" button to reload the list if you add new files while the application is running.

### Browser Integration

When grid generation completes successfully, the HTML file automatically opens in your default browser, allowing you to immediately view the generated hex grid with the item legend and scrolling functionality.

### Automatic Naming

Output files are automatically named using the pattern:
`{yaml-filename}-{YYYY-MM-DD-HH-MM-SS}`

For example:
- Input: `fantasy-world.yaml`
- Output: `fantasy-world-2024-01-15-14-30-25.svg` and `fantasy-world-2024-01-15-14-30-25.html`

## YAML Configuration Format

Create a YAML file with the following structure:

```yaml
default: "#F5F5DC"
items:
  - name: "Forest"
    percentage: 30.0
    style: "fill"
    color: "#228B22"
  
  - name: "Water"
    percentage: 15.0
    style: "fill"
    color: "#4169E1"
  
  - name: "Village"
    percentage: 10.0
    style: "dot"
    color: "#FFD700"
```

### Configuration Options

- **default**: Hex color code for the background color of empty cells and dot-style items
- **name**: A descriptive name for the item type
- **percentage**: Percentage of grid cells to fill with this item (0-100)
- **style**: Either "fill" (colored hexagon) or "dot" (colored dot in center)
- **color**: Hex color code (e.g., "#FF0000" for red)

### Rules

- **default** color is required
- Total percentage should not exceed 100%
- Valid styles are "fill" and "dot"
- Use valid hex color codes

## Output Files

The application generates two files:

1. **SVG file** (`.svg`): Vector graphics file containing the hex grid
2. **HTML file** (`.html`): Web page with embedded SVG, scrolling, and item legend

## Hex Grid Layout

The hexagons are arranged in a proper staggered pattern where:
- Each hexagon touches its neighbors without overlapping
- Odd rows are offset by half a hexagon width
- The grid forms a true hexagonal tiling pattern
- Each hexagon has 6 sides that can connect to adjacent hexagons

## Example

A sample configuration file `sample_config.yaml` is included with the following items:
- Forest (30% - green fill)
- Water (15% - blue fill)
- Mountain (20% - brown fill)
- Village (10% - yellow dot)
- Road (5% - gray dot)
- Empty (20% - beige fill)

## Building

To create an executable:

```bash
go build -o hexgrid .
```

Then run the executable:

```bash
./hexgrid
```

## Dependencies

- **Fyne v2**: Cross-platform GUI framework
- **YAML v3**: YAML parsing library

## License

This project is open source and available under the MIT License.
