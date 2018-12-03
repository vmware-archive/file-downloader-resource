package file

import (
	"path"
	"path/filepath"
	"strings"
)

func Matches(fileFullName, productSlug, pattern, version string) (bool, error) {
	matchesPattern, err := filepath.Match(path.Join(productSlug, pattern), fileFullName)
	if err != nil {
		return false, err
	}
	if matchesPattern {
		return strings.Contains(fileFullName, version), nil
	}
	return false, nil
}
