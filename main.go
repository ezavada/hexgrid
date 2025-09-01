package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	YAMLPath   string
	OutputPath string
	GridRows   int
	GridCols   int
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hex Grid Generator")
	myWindow.Resize(fyne.NewSize(600, 400))

	config := &Config{
		GridRows: 10,
		GridCols: 10,
	}

	// Output path label (declared early so it can be used in YAML selection)
	outputPathLabel := widget.NewLabel("No output file selected")

	// YAML file selection dropdown
	yamlSelectLabel := widget.NewLabel("Select YAML Configuration:")
	yamlSelectDropdown := widget.NewSelect([]string{}, func(selected string) {
		if selected != "" {
			// Get current working directory
			currentDir, err := os.Getwd()
			if err != nil {
				currentDir = "."
			}
			gridSpecsDir := filepath.Join(currentDir, "grid-specs")

			// Set the full path to the selected file
			config.YAMLPath = filepath.Join(gridSpecsDir, selected)

			// Auto-generate output path when YAML file is selected
			generateOutputPath(config, outputPathLabel)
		}
	})

	// Function to populate the dropdown with YAML files
	populateYAMLDropdown := func() {
		// Get the executable path to find the app bundle location
		execPath, err := os.Executable()
		if err != nil {
			execPath = "."
		}

		// Check if we're running from an app bundle
		appBundlePath := filepath.Join(filepath.Dir(execPath), "..", "..", "..")
		gridSpecsDir := filepath.Join(appBundlePath, "Contents", "Resources", "grid-specs")

		// If not in app bundle, try current directory
		if _, err := os.Stat(gridSpecsDir); os.IsNotExist(err) {
			currentDir, err := os.Getwd()
			if err != nil {
				currentDir = "."
			}
			gridSpecsDir = filepath.Join(currentDir, "grid-specs")
		}

		// Create grid-specs directory if it doesn't exist
		if _, err := os.Stat(gridSpecsDir); os.IsNotExist(err) {
			os.MkdirAll(gridSpecsDir, 0755)
		}

		// Read all YAML files from the grid-specs directory
		files, err := os.ReadDir(gridSpecsDir)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to read grid-specs directory: %w", err), myWindow)
			return
		}

		var yamlFiles []string
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
				yamlFiles = append(yamlFiles, file.Name())
			}
		}

		if len(yamlFiles) == 0 {
			yamlFiles = []string{"No YAML files found"}
		}

		yamlSelectDropdown.SetOptions(yamlFiles)
		yamlSelectDropdown.SetSelected("")
	}

	// Refresh button for the dropdown
	refreshBtn := widget.NewButton("Refresh", populateYAMLDropdown)

	// Populate the dropdown on startup
	populateYAMLDropdown()

	// Grid size inputs
	rowsInput := widget.NewEntry()
	rowsInput.SetText("10")
	rowsInput.OnChanged = func(value string) {
		if val, err := parseInt(value); err == nil {
			config.GridRows = val
		}
	}

	colsInput := widget.NewEntry()
	colsInput.SetText("10")
	colsInput.OnChanged = func(value string) {
		if val, err := parseInt(value); err == nil {
			config.GridCols = val
		}
	}

	// Output file selection
	outputSelectBtn := widget.NewButton("Select Output File", func() {
		// Get the executable path to find the app bundle location
		execPath, err := os.Executable()
		if err != nil {
			execPath = "."
		}

		// Check if we're running from an app bundle
		appBundlePath := filepath.Join(filepath.Dir(execPath), "..", "..", "..")
		generatedGridsDir := filepath.Join(appBundlePath, "Contents", "Resources", "generated-grids")

		// If not in app bundle, try current directory
		if _, err := os.Stat(generatedGridsDir); os.IsNotExist(err) {
			currentDir, err := os.Getwd()
			if err != nil {
				currentDir = "."
			}
			generatedGridsDir = filepath.Join(currentDir, "generated-grids")
		}

		// Create generated-grids directory if it doesn't exist
		if _, err := os.Stat(generatedGridsDir); os.IsNotExist(err) {
			os.MkdirAll(generatedGridsDir, 0755)
		}

		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()

			config.OutputPath = writer.URI().Path()
			outputPathLabel.SetText(filepath.Base(config.OutputPath))
		}, myWindow)
	})

	// Generate button
	generateBtn := widget.NewButton("Generate Hex Grid", func() {
		if config.YAMLPath == "" {
			dialog.ShowError(fmt.Errorf("please select a YAML file"), myWindow)
			return
		}
		if config.OutputPath == "" {
			dialog.ShowError(fmt.Errorf("please select an output file"), myWindow)
			return
		}

		err := generateHexGrid(config)
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			// Open the generated HTML file in the browser
			htmlPath := config.OutputPath + ".html"
			openInBrowser(htmlPath)
			dialog.ShowInformation("Success", "Hex grid generated successfully and opened in browser!", myWindow)
		}
	})

	// Layout
	form := container.NewVBox(
		widget.NewLabel("Hex Grid Generator"),
		widget.NewSeparator(),
		yamlSelectLabel,
		container.NewHBox(yamlSelectDropdown, refreshBtn),
		widget.NewSeparator(),
		container.NewHBox(
			container.NewVBox(
				widget.NewLabel("Grid Rows:"),
				rowsInput,
			),
			container.NewVBox(
				widget.NewLabel("Grid Columns:"),
				colsInput,
			),
		),
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel("Output File:"), outputSelectBtn),
		outputPathLabel,
		widget.NewSeparator(),
		generateBtn,
	)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// generateOutputPath automatically generates an output path based on the YAML file name and timestamp
func generateOutputPath(config *Config, outputPathLabel *widget.Label) {
	if config.YAMLPath == "" {
		return
	}

	// Get the base name of the YAML file (without extension)
	yamlBaseName := filepath.Base(config.YAMLPath)
	yamlNameWithoutExt := yamlBaseName[:len(yamlBaseName)-len(filepath.Ext(yamlBaseName))]

	// Generate timestamp
	timestamp := time.Now().Format("2006-01-02-15-04-05")

	// Get the executable path to find the app bundle location
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}

	// Check if we're running from an app bundle
	appBundlePath := filepath.Join(filepath.Dir(execPath), "..", "..", "..")
	generatedGridsDir := filepath.Join(appBundlePath, "Contents", "Resources", "generated-grids")

	// If not in app bundle, try current directory
	if _, err := os.Stat(generatedGridsDir); os.IsNotExist(err) {
		currentDir, err := os.Getwd()
		if err != nil {
			currentDir = "."
		}
		generatedGridsDir = filepath.Join(currentDir, "generated-grids")
	}

	// Create generated-grids directory if it doesn't exist
	if _, err := os.Stat(generatedGridsDir); os.IsNotExist(err) {
		os.MkdirAll(generatedGridsDir, 0755)
	}

	// Generate output path
	outputFileName := fmt.Sprintf("%s-%s", yamlNameWithoutExt, timestamp)
	config.OutputPath = filepath.Join(generatedGridsDir, outputFileName)

	// Update the label
	outputPathLabel.SetText(outputFileName)
}

// openInBrowser opens the specified file in the default browser
func openInBrowser(filePath string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // Linux and other Unix-like systems
		cmd = exec.Command("xdg-open", filePath)
	}

	// Run the command in the background
	go func() {
		err := cmd.Run()
		if err != nil {
			// Silently fail - browser opening is not critical
			fmt.Printf("Failed to open browser: %v\n", err)
		}
	}()
}

func generateHexGrid(config *Config) error {
	// Load YAML configuration
	yamlConfig, err := LoadYAMLConfig(config.YAMLPath)
	if err != nil {
		return fmt.Errorf("failed to load YAML config: %w", err)
	}

	// Create hex grid
	grid := CreateHexGrid(config.GridRows, config.GridCols, yamlConfig)

	// Populate grid with items
	grid.PopulateGrid()

	// Generate SVG file
	svgPath := config.OutputPath + ".svg"
	err = GenerateSVG(grid, svgPath)
	if err != nil {
		return fmt.Errorf("failed to generate SVG: %w", err)
	}

	// Generate HTML file
	htmlPath := config.OutputPath + ".html"
	err = GenerateHTML(grid, svgPath, htmlPath)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	return nil
}
