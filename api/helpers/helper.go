package helpers

import (
	"os"
	"strings"
)

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") { // 与环境变量中的域名相同
		return false
	}

	// 去除http://(https://)、www.
	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]

	if newUrl == os.Getenv("DOMAIN") {
		return false
	}

	return true
}

func ForceHTTPS(url string) string {
	if url[:4] != "http" {
		return "https://" + url
	}
	return url
}
