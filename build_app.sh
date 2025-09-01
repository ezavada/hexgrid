#!/bin/bash

# Build the Go application
echo "Building Go application..."
go build -o hexgrid .

# Create app bundle structure
echo "Creating app bundle..."
mkdir -p HexGrid.app/Contents/MacOS
mkdir -p HexGrid.app/Contents/Resources

# Copy executable
cp hexgrid HexGrid.app/Contents/MacOS/

# Make executable
chmod +x HexGrid.app/Contents/MacOS/hexgrid

# Copy grid-specs folder to Resources
cp -r grid-specs HexGrid.app/Contents/Resources/

# Create generated-grids folder in Resources
mkdir -p HexGrid.app/Contents/Resources/generated-grids

# Set proper permissions
chmod -R 755 HexGrid.app

echo "App bundle created: HexGrid.app"
echo "You can now double-click HexGrid.app to run the application"
