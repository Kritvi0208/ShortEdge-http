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


