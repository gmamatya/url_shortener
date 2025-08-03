package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return strings.Replace(url, "http://", "https://", 1)
}

func IsServiceDomain(url string) bool {
	domain := os.Getenv("DOMAIN")
	return strings.HasPrefix(url, "http://"+domain) || strings.HasPrefix(url, "https://"+domain)
}
