package search

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func IP(ip, endpoint string) (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(endpoint)
	// Always return an error if we get one.
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// That means we found an IP address.
		// NOTE: THIS IS BAD.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		name := gjson.Get(string(body), "name").String()
		return ip, fmt.Errorf("found IP address %s in project %s", ip, name)
	}
	return "", nil
}
