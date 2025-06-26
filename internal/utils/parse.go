package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/avct/uasurfer"
)

type Location struct {
	Country string
	City    string
}

func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func GetLocation(ip string) (Location, error) {
	var loc Location
	url := fmt.Sprintf("https://ipwho.is/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return loc, err
	}
	defer resp.Body.Close()

	var res struct {
		Success bool   `json:"success"`
		Country string `json:"country"`
		City    string `json:"city"`
	}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil || !res.Success {
		return loc, errors.New("failed to get location")
	}

	loc.Country = res.Country
	loc.City = res.City
	return loc, nil
}

func ParseUserAgent(ua string) (browser, os, device string) {
	info := uasurfer.Parse(ua)
	browser = trimPrefix(info.Browser.Name.String(), "Browser")
	device = trimPrefix(info.DeviceType.String(), "Device")
	os = info.OS.Name.String()
	return
}


func trimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

// package utils

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"

// 	"github.com/avct/uasurfer"
// )

// // GetCountryFromIP uses ipwho.is API to get the country from IP
// func GetCountryFromIP(ip string) string {
// 	resp, err := http.Get(fmt.Sprintf("https://ipwho.is/%s", ip))
// 	if err != nil {
// 		return "Unknown"
// 	}
// 	defer resp.Body.Close()

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	var data struct {
// 		Country string `json:"country"`
// 	}
// 	_ = json.Unmarshal(body, &data)

// 	if data.Country == "" {
// 		return "Unknown"
// 	}
// 	return data.Country
// }

// // GetDeviceFromUserAgent parses user-agent string
// func GetDeviceFromUserAgent(ua string) string {
// 	uaInfo := uasurfer.Parse(ua)
// 	return strings.TrimPrefix(uaInfo.DeviceType.String(), "Device")

// }
