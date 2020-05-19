package helpers

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	requestTimeOut = 5
)

// ServerSideParameters :
type ServerSideParameters struct {
	URL         string     `json:"URL"`
	URLQuery    url.Values `json:"Queries"`
	UserAgent   string     `json:"User-Agent"`
	XForwardFor string     `json:"X-Fowarded-For"`
	Referer     string     `json:"Referer"`
}

// ServerSideCall :
func ServerSideCall(uri string, ssp ServerSideParameters) ([]byte, error) {
	timeout := time.Duration(time.Duration(requestTimeOut) * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return []byte(``), err
	}

	request.Header.Set("User-Agent", ssp.UserAgent)
	request.Header.Set("X-Forwarded-For", ssp.XForwardFor)

	resp, err := client.Do(request)
	if err != nil {
		return []byte(``), err
	}

	defer resp.Body.Close()
	contentResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(``), err
	}

	return contentResponse, nil
}

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
