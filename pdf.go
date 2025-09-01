package main

import (
	"fmt"
	"math"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF creates a PDF representation of the hex grid
func GeneratePDF(grid *HexGrid, outputPath string) error {
	// Create new PDF document
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 10)

	// Calculate page dimensions in mm
	pageWidth, pageHeight := pdf.GetPageSize()
	margin := 20.0
	usableWidth := pageWidth - 2*margin
	usableHeight := pageHeight - 2*margin

	// Calculate hexagon size for PDF (convert from SVG units to mm)
	// SVG uses pixels, PDF uses mm - approximate conversion
	hexSizeMM := 8.0                // Size of hexagon in mm
	hexWidthMM := hexSizeMM * 1.732 // Width (2 * size * cos(30°))
	hexHeightMM := hexSizeMM * 2.0  // Height

	// Calculate grid dimensions
	gridWidth := float64(grid.Cols) * hexWidthMM
	gridHeight := float64(grid.Rows) * (hexHeightMM / 2)

	// Center the grid on the page
	startX := margin + (usableWidth-gridWidth)/2
	startY := margin + (usableHeight-gridHeight)/2

	// Draw hexagons
	for row := 0; row < grid.Rows; row++ {
		for col := 0; col < grid.Cols; col++ {
			cell := grid.Cells[row][col]

			// Calculate hexagon center position
			x := startX + float64(col)*hexWidthMM
			if row%2 == 1 {
				x += hexWidthMM / 2
			}
			y := startY + float64(row)*(hexHeightMM/2)

			// Draw hexagon
			drawHexagon(pdf, x, y, hexSizeMM, cell, grid.DefaultColor)

			// Add dice result if available
			if cell.DiceResult != nil {
				// Position text to the right of the hexagon
				textX := x + hexSizeMM + 2
				textY := y + 1

				pdf.SetFont("Arial", "", 8)
				pdf.SetTextColor(0, 0, 0)
				pdf.Text(textX, textY, fmt.Sprintf("%d", cell.DiceResult.Total))
			}
		}
	}

	// Add legend
	addLegend(pdf, grid, pageWidth, pageHeight, margin)

	// Save PDF
	return pdf.OutputFileAndClose(outputPath)
}

// drawHexagon draws a hexagon at the specified position
func drawHexagon(pdf *gofpdf.Fpdf, centerX, centerY, size float64, cell *HexCell, defaultColor string) {
	// Calculate hexagon points
	var points []gofpdf.PointType
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3.0
		x := centerX + size*math.Cos(angle)
		y := centerY + size*math.Sin(angle)
		points = append(points, gofpdf.PointType{X: x, Y: y})
	}

	// Determine fill color
	var fillColor string
	var strokeColor string

	if cell.ItemType != nil {
		if cell.ItemType.Style == "fill" {
			fillColor = cell.ItemType.Color
		} else {
			fillColor = defaultColor
		}
		strokeColor = "#333333"
	} else {
		fillColor = defaultColor
		strokeColor = "#CCCCCC"
	}

	// Convert hex color to RGB
	r, g, b := hexToRGB(fillColor)
	pdf.SetFillColor(r, g, b)

	r, g, b = hexToRGB(strokeColor)
	pdf.SetDrawColor(r, g, b)

	// Draw hexagon
	pdf.Polygon(points, "F") // Fill
	pdf.Polygon(points, "D") // Draw outline

	// Add dot if style is "dot"
	if cell.ItemType != nil && cell.ItemType.Style == "dot" {
		r, g, b = hexToRGB(cell.ItemType.Color)
		pdf.SetFillColor(r, g, b)
		pdf.Circle(centerX, centerY, 2, "F")
	}
}

// addLegend adds a legend to the PDF
func addLegend(pdf *gofpdf.Fpdf, grid *HexGrid, pageWidth, pageHeight, margin float64) {
	// Position legend in top-right corner
	legendX := pageWidth - margin - 60
	legendY := margin

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(legendX, legendY, "Item Legend")

	pdf.SetFont("Arial", "", 10)
	yOffset := legendY + 8

	for _, itemType := range grid.ItemTypes {
		// Draw symbol
		symbolX := legendX
		symbolY := yOffset - 3

		if itemType.Style == "fill" {
			// Draw filled square
			r, g, b := hexToRGB(itemType.Color)
			pdf.SetFillColor(r, g, b)
			pdf.Rect(symbolX, symbolY, 4, 4, "F")
		} else {
			// Draw circle with dot
			pdf.SetFillColor(255, 255, 255)
			pdf.Circle(symbolX+2, symbolY+2, 2, "F")
			pdf.SetDrawColor(0, 0, 0)
			pdf.Circle(symbolX+2, symbolY+2, 2, "D")

			r, g, b := hexToRGB(itemType.Color)
			pdf.SetFillColor(r, g, b)
			pdf.Circle(symbolX+2, symbolY+2, 1, "F")
		}

		// Draw text
		pdf.SetTextColor(0, 0, 0)
		text := fmt.Sprintf("%s (%.1f%%)", itemType.Name, itemType.Percentage)
		if itemType.Dice != "" {
			text += fmt.Sprintf(" - %s", itemType.Dice)
		}
		pdf.Text(symbolX+8, symbolY+3, text)

		yOffset += 6
	}
}

// hexToRGB converts hex color string to RGB values
func hexToRGB(hex string) (int, int, int) {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Default to black if invalid
	if len(hex) != 6 {
		return 0, 0, 0
	}

	// Parse RGB values
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}
