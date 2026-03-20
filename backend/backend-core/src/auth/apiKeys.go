package auth

import (
	"net"
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func extractAPIKey(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-API-Key"))
}

func hashAPIKey(key string) string {
	return sharedUtils.GenerateHexHash(key)
}

func isIPAllowed(ip string, restrictions []dbModel.APIKeyIPRestrictionEntity) bool {
	if len(restrictions) == 0 {
		return true
	}

	parsedIP := net.ParseIP(strings.TrimSpace(ip))
	if parsedIP == nil {
		return false
	}

	for _, restriction := range restrictions {
		_, cidrNet, err := net.ParseCIDR(strings.TrimSpace(restriction.CIDR))
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
