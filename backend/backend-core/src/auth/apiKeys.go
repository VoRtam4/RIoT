/**
 * @file apiKeys.go
 * @brief Pomocné funkce pro extrakci API klíčů a vyhodnocení jejich IP omezení.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package auth

import (
	"net"
	"net/http"
	"strings"
)

func extractAPIKey(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-API-Key"))
}

func isIPAllowed(ip string, restrictions []string) bool {
	if len(restrictions) == 0 {
		return true
	}
	parsedIP := net.ParseIP(strings.TrimSpace(ip))
	if parsedIP == nil {
		return false
	}
	for _, restriction := range restrictions {
		_, cidrNet, err := net.ParseCIDR(restriction)
		if err != nil {
			continue
		}
		if cidrNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}

func extractClientIP(r *http.Request) string {
	ip := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
