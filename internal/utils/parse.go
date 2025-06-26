package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/avct/uasurfer"
)

// GetCountryFromIP uses ipwho.is API to get the country from IP
func GetCountryFromIP(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://ipwho.is/%s", ip))
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		Country string `json:"country"`
	}
	_ = json.Unmarshal(body, &data)

	if data.Country == "" {
		return "Unknown"
	}
	return data.Country
}

// GetDeviceFromUserAgent parses user-agent string
func GetDeviceFromUserAgent(ua string) string {
	uaInfo := uasurfer.Parse(ua)
	return strings.TrimPrefix(uaInfo.DeviceType.String(), "Device")

}
