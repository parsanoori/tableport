package tableport

import "unicode"

// isPersian checks if a rune belongs to the Persian character set
func isPersian(r rune) bool {
	// Persian characters in Unicode ranges
	return (r >= 0x0600 && r <= 0x06FF) || // Arabic (including Persian) block
		(r >= 0x0750 && r <= 0x077F) || // Arabic Supplement
		(r >= 0x08A0 && r <= 0x08FF) || // Arabic Extended-A
		(r >= 0xFB50 && r <= 0xFDFF) || // Arabic Presentation Forms-A
		(r >= 0xFE70 && r <= 0xFEFF) || // Arabic Presentation Forms-B
		(r >= 0x10E60 && r <= 0x10E7F) // Rumi Numeral Symbols (used in Persian)
}

// isAllPersian checks if all characters in a string are Persian
func isAllPersian(s string) bool {
	for _, r := range s {
		if !isPersian(r) && !unicode.IsSpace(r) { // Allow spaces
			return false
		}
	}
	return true
}
