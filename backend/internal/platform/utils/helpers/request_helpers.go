package helpers

import (
	"fmt"
	"net"
	"strings"
)

// Helpers (moved from old integration file)
func ParseUserAgent(userAgent string) string {
	if userAgent == "" {
		return "Unknown device"
	}
	ua := strings.ToLower(userAgent)
	var os string
	if strings.Contains(ua, "windows") {
		os = "Windows"
	} else if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		os = "macOS"
	} else if strings.Contains(ua, "linux") {
		os = "Linux"
	} else if strings.Contains(ua, "android") {
		os = "Android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "iOS"
	} else {
		os = "Unknown OS"
	}
	var browser string
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "edge") {
		browser = "Chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") {
		browser = "Safari"
	} else if strings.Contains(ua, "edge") {
		browser = "Edge"
	} else if strings.Contains(ua, "opera") {
		browser = "Opera"
	} else {
		browser = "Unknown browser"
	}
	return fmt.Sprintf("%s on %s", browser, os)
}

func GetLocationFromIP(ip string) string {
	if ip == "" {
		return "Unknown location"
	}
	if net.ParseIP(ip) == nil {
		return "Unknown location"
	}
	if IsPrivateIP(ip) {
		return "Local network"
	}
	return fmt.Sprintf("IP: %s", ip)
}

func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 127 {
			return true
		}
		if ip4[0] == 10 {
			return true
		}
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	}
	if ip.IsLoopback() {
		return true
	}
	return false
}
