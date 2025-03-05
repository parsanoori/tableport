package tableport

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed assets/fonts/Vazir.ttf
var myFont []byte

var tempFontPath string // Store the temp font path

// init() runs when the package is imported
func init() {
	tempDir := os.TempDir() // Get OS temp directory
	tempFontPath = filepath.Join(tempDir, "embedded_font.ttf")

	// Write font to the temp file
	err := os.WriteFile(tempFontPath, myFont, 0644)
	if err != nil {
		panic("Failed to write the embedded font to a temp file")
	}

	// Optional: Delete the temp file when the program exits
	go func() {
		<-make(chan struct{}) // Prevent premature cleanup (adjust if needed)
		os.Remove(tempFontPath)
	}()
}
