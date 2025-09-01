package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// Hexagon parameters
const (
	HexSize   = 25.0 // Distance from center to any corner
	HexWidth  = 43.3 // Width of hexagon (2 * HexSize * cos(30°))
	HexHeight = 50.0 // Height of hexagon (2 * HexSize)
	// HexSpacing = 0.0  // Space between hexagons (0 for touching)
	HexColumnOffset       = 40.0                                           // Offset for the x axis between columns
	HexEvenRowStartOffset = (HexWidth-HexColumnOffset)/2 + HexColumnOffset // Offset for the start of the even rows
)

// GenerateSVG creates an SVG representation of the hex grid
func GenerateSVG(grid *HexGrid, outputPath string) error {
	// Calculate SVG dimensions for proper hex grid layout
	svgWidth := float64(grid.Cols) * (HexWidth + HexColumnOffset)
	if grid.Rows > 1 {
		svgWidth += HexEvenRowStartOffset // add the width of the offset for the even rows
	}
	// Height: each row takes HexHeight/2 (half the cell height)
	svgHeight := float64(grid.Rows) * (HexHeight / 2)
	if grid.Rows%2 == 0 {
		// even number of rows, so the last row is offset by half the cell height
		svgHeight += HexHeight / 2
	}

	// Start SVG content
	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="%.1f" height="%.1f" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style>
      .hexagon { stroke: #333; stroke-width: 1; }
      .hexagon-dot { fill: none; }
    </style>
  </defs>
  <g transform="translate(20, 20)">`, svgWidth, svgHeight)

	// Generate hexagons
	for row := 0; row < grid.Rows; row++ {
		for col := 0; col < grid.Cols; col++ {
			cell := grid.Cells[row][col]

			// Calculate hexagon center position for proper staggered layout
			// Each hexagon is offset by half its width in odd rows
			x := float64(col) * (HexWidth + HexColumnOffset)
			if row%2 == 1 {
				x += HexEvenRowStartOffset
			}

			// Each row is offset by half the cell height
			y := float64(row) * (HexHeight / 2)

			cell.X = x
			cell.Y = y

			// Generate hexagon path
			hexPath := generateHexagonPath(x, y)

			// Determine styling based on item type
			var fillColor, strokeColor string

			if cell.ItemType != nil {
				if cell.ItemType.Style == "fill" {
					// For fill style, use the item's color
					fillColor = cell.ItemType.Color
				} else {
					// For dot style, use the default background color
					fillColor = grid.DefaultColor
				}
				strokeColor = "#333"
			} else {
				// For empty cells, use the default color
				fillColor = grid.DefaultColor
				strokeColor = "#ccc"
			}

			// Add hexagon with direct color attributes
			svg += fmt.Sprintf(`
    <path d="%s" fill="%s" stroke="%s" stroke-width="1"/>`, hexPath, fillColor, strokeColor)

			// Add dot if style is "dot" (3x bigger with black outline)
			if cell.ItemType != nil && cell.ItemType.Style == "dot" {
				svg += fmt.Sprintf(`
    <circle cx="%.1f" cy="%.1f" r="9" fill="%s" stroke="black" stroke-width="2"/>`, x, y, cell.ItemType.Color)
			}
		}
	}

	svg += `
  </g>
</svg>`

	// Write SVG to file
	err := os.WriteFile(outputPath, []byte(svg), 0644)
	if err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	return nil
}

// generateHexagonPath creates the SVG path for a hexagon
func generateHexagonPath(centerX, centerY float64) string {
	var points []string

	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3.0
		x := centerX + HexSize*math.Cos(angle)
		y := centerY + HexSize*math.Sin(angle)

		if i == 0 {
			points = append(points, fmt.Sprintf("M %.1f %.1f", x, y))
		} else {
			points = append(points, fmt.Sprintf("L %.1f %.1f", x, y))
		}
	}

	points = append(points, "Z")
	return strings.Join(points, " ")
}

// GenerateHTML creates an HTML page that embeds the SVG with scrolling and legend
func GenerateHTML(grid *HexGrid, svgPath, outputPath string) error {
	// Read the SVG content
	svgContent, err := os.ReadFile(svgPath)
	if err != nil {
		return fmt.Errorf("failed to read SVG file: %w", err)
	}

	// Create legend HTML
	legendHTML := `<div class="legend">
    <h3>Item Legend</h3>
    <div class="legend-items">`

	for _, itemType := range grid.ItemTypes {
		var symbol string
		if itemType.Style == "fill" {
			symbol = fmt.Sprintf(`<div class="legend-symbol fill" style="background-color: %s;"></div>`, itemType.Color)
		} else {
			symbol = fmt.Sprintf(`<div class="legend-symbol dot"><div class="dot" style="background-color: %s;"></div></div>`, itemType.Color)
		}

		legendHTML += fmt.Sprintf(`
      <div class="legend-item">
        %s
        <span class="legend-name">%s (%.1f%%)</span>
      </div>`, symbol, itemType.Name, itemType.Percentage)
	}

	legendHTML += `
    </div>
  </div>`

	// Create HTML content
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Hex Grid Generator</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: Arial, sans-serif;
            background-color: #f5f5f5;
        }
        .container {
            display: flex;
            gap: 20px;
            max-width: 100%%;
        }
        .svg-container {
            flex: 1;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: auto;
            max-height: 80vh;
        }
        .svg-container svg {
            display: block;
            margin: 0;
        }
        .legend {
            width: 250px;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            height: fit-content;
        }
        .legend h3 {
            margin-top: 0;
            color: #333;
        }
        .legend-items {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        .legend-item {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .legend-symbol {
            width: 20px;
            height: 20px;
            border: 1px solid #333;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .legend-symbol.fill {
            border-radius: 0;
        }
        .legend-symbol.dot {
            border-radius: 50%%;
            background: white;
        }
        .legend-symbol .dot {
            width: 8px;
            height: 8px;
            border-radius: 50%%;
        }
        .legend-name {
            font-size: 14px;
            color: #555;
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <h1>Hex Grid Generator</h1>
    <div class="container">
        <div class="svg-container">
            %s
        </div>
        %s
    </div>
</body>
</html>`, string(svgContent), legendHTML)

	// Write HTML to file
	err = os.WriteFile(outputPath, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}
