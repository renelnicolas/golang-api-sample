package helpers

import (
	"net"
	"net/http"
	"strings"
)

// GetIP defines user ip & proxy services
func GetIP(remoteAddr string, headers http.Header) string {
	if clientIP, _, err := net.SplitHostPort(remoteAddr); nil == err && clientIP != "::1" && clientIP != "127.0.0.1" {
		if prior, ok := headers["X-Forwarded-For"]; ok {
			ipSet := map[string]bool{clientIP: true}

			for _, ip := range strings.Split(strings.Join(prior, ","), ",") {
				ip = strings.TrimFunc(ip, TrimWhitespaceFn)
				if !ipSet[ip] {
					ipSet[ip] = true

					clientIP = clientIP + ", " + strings.TrimFunc(ip, TrimWhitespaceFn)
				}
			}
		}

		return clientIP
	}

	return ""
}
