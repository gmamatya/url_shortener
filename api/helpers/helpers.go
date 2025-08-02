package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	// Ensure the URL uses HTTPS
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func RemoveDomainError(url string) bool {
	// Check if the URL's domain is valid
	if url == os.Getenv("DOMAIN") {
		return false // Invalid domain
	}

	// Remove the scheme and www prefix for domain validation
	newURL := strings.Replace(url, "http://", "", 1)    // Remove http scheme
	newURL = strings.Replace(newURL, "https://", "", 1) // Remove https scheme
	newURL = strings.Replace(newURL, "www.", "", 1)     // Remove www prefix
	newURL = strings.Split(newURL, "/")[0]              // Get the domain part

	if newURL == os.Getenv("DOMAIN") {
		return false // Invalid domain
	}

	return true // Valid domain
}
