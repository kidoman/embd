package embd

import "path/filepath"

// FindFirstMatchingFile finds the first glob match in the filesystem.
// Inspiration: https://github.com/mrmorphic/hwio/blob/master/hwio.go#L451
func FindFirstMatchingFile(glob string) (string, error) {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}
	if len(matches) >= 1 {
		return matches[0], nil
	}
	return "", nil
}
