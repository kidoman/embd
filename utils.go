package embd

import "path/filepath"

// Inspiration: https://github.com/mrmorphic/hwio/blob/master/hwio.go#L451
func findFirstMatchingFile(glob string) (string, error) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}
	if len(matches) >= 1 {
		return matches[0], nil
	}
	return "", nil
}
