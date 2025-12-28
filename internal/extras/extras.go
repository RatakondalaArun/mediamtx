package extras

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bluenviron/mediamtx/internal/conf"
)

func LoadConfigFromAPI(apiURL string) (*conf.Conf, error) {
	if apiURL == "" {
		return nil, fmt.Errorf("API URL cannot be empty")
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}
	var c conf.Conf
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	return &c, nil
}
