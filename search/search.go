package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/darron/gips/core"
)

func IP(ip, endpoint string) (core.Project, error) {
	var project core.Project
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(endpoint)
	// Always return an error if we get one.
	if err != nil {
		return project, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// That means we found an IP address.
		// NOTE: THIS IS BAD.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return project, err
		}
		err = json.Unmarshal(body, &project)
		if err != nil {
			return project, fmt.Errorf("error unmarshalling JSON: %s", err)
		}
		return project, fmt.Errorf("found IP address %s in project %s", ip, project.Name)
	}
	return project, nil
}
