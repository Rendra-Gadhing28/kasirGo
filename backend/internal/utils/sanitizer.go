package utils

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var strictPolicy = bluemonday.StrictPolicy()
var nonSafeFilenameRe = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

// SanitizeString strips any HTML or scripting tags
func SanitizeString(s string) string {
	cleaned := strictPolicy.Sanitize(s)
	return strings.TrimSpace(cleaned)
}

// SanitizeFilename prevents path traversal and unsafe characters
func SanitizeFilename(name string) string {
	base := filepath.Base(name)
	return nonSafeFilenameRe.ReplaceAllString(base, "_")
}
