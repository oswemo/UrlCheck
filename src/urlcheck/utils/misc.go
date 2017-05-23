package utils

// Miscellaneous utility functions

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// URLDecode attempts to decode a string using standard URL decoding libs.
// The purpose of this method is primarily to wrap a log statement without having
// to redo everytime I want to decode something.
func URLDecode(content string) (string, error) {
	decoded, err := url.PathUnescape(content)
	if err != nil {
		LogError(LogFields{}, err, "Failed to decode query path")
		return "", err
	}
	return decoded, nil
}

// validateUrl attempts to validate that the host and port combination are valid.
// That is, the
func ValidateUrl(hostname string) error {
	// A hostname should contain the host and port.
	parts := strings.Split(hostname, ":")
	if len(parts) != 2 {
		return errors.New("Invalid hostname port combination given")
	}

	// Hostname matching regex.
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if !re.MatchString(parts[0]) {
		return errors.New("Hostname does not appear to be a valid format")
	}

	// Check the port range.
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("Port does not appear to be valid")
	}

	if port < 1 || port > 65535 {
		return errors.New("Port is not within the valid range")
	}

	return nil
}
