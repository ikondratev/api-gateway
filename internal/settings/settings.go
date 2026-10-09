package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	settingsPath = "settings"
	settingsName = "settings.json"
)

type Settings struct {
	Environment string   `json:"env"`
	Server      Server   `json:"server"`
	Upstream    Upstream `json:"upstream"`
}

type Upstream struct {
	EventSerivce EventService `json:"event_service"`
}

type EventService struct {
	Url string `json:"url"`
}

type Server struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	WaitingShutdown int    `json:"waiting_shutdown"`
	HeaderTimeout   int    `json:"header_timeout"`
	ReadTimeout     int    `json:"read_timeout"`
	WriteTimeout    int    `json:"write_teimeout"`
	IdleTimeout     int    `json:"idle_timeout"`
}

func New(env string) (*Settings, error) {
	path := fmt.Sprintf("%s/%s_%s", settingsPath, env, settingsName)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load file: %w", err)
	}

	var set *Settings
	if err := json.Unmarshal(data, &set); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	if err := validate(set); err != nil {
		return nil, err
	}

	return set, nil
}

func validate(s *Settings) error {
	if s.Upstream.EventSerivce.Url == "" {
		return fmt.Errorf("upstream for event service is empty")
	}

	return nil
}
