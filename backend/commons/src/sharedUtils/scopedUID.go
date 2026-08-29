package sharedUtils

import (
	"fmt"
	"strings"
)

func NormalizePrefixedUID(inputUID string, prefix string, entityName string) (string, string, error) {
	trimmedPrefix := strings.TrimSpace(prefix)
	if trimmedPrefix == "" {
		return "", "", fmt.Errorf("%s UID prefix must not be empty", entityName)
	}
	suffix := strings.TrimSpace(inputUID)
	canonicalPrefix := trimmedPrefix + ":"
	if strings.HasPrefix(suffix, canonicalPrefix) {
		suffix = strings.TrimPrefix(suffix, canonicalPrefix)
	}
	suffix = strings.Join(strings.Fields(suffix), "_")
	if suffix == "" {
		return "", "", fmt.Errorf("%s UID suffix must not be empty", entityName)
	}
	if strings.ContainsAny(suffix, ".:") {
		return "", "", fmt.Errorf("%s UID suffix must not contain '.' or ':'", entityName)
	}
	return canonicalPrefix + suffix, suffix, nil
}

func NormalizeScopedUID(inputUID string, scopeUID string, childPrefix string, entityName string) (string, string, error) {
	trimmedScopeUID := strings.TrimSpace(scopeUID)
	if trimmedScopeUID == "" {
		return "", "", fmt.Errorf("%s scope UID must not be empty", entityName)
	}
	trimmedChildPrefix := strings.TrimSpace(childPrefix)
	if trimmedChildPrefix == "" {
		return "", "", fmt.Errorf("%s child UID prefix must not be empty", entityName)
	}
	prefix := trimmedScopeUID + "." + trimmedChildPrefix + ":"
	suffix := strings.TrimSpace(inputUID)
	if strings.HasPrefix(suffix, prefix) {
		suffix = strings.TrimPrefix(suffix, prefix)
	} else if strings.Contains(suffix, ".") {
		return "", "", fmt.Errorf("%s UID must be in scope %s", entityName, trimmedScopeUID)
	}
	uid, normalizedSuffix, err := NormalizePrefixedUID(suffix, trimmedChildPrefix, entityName)
	if err != nil {
		return "", "", err
	}
	return trimmedScopeUID + "." + uid, normalizedSuffix, nil
}
