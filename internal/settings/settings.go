package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	settingsPath = "settings"
	settingsName = "settings.json"
)

type Settings struct {
	Environment string   `json:"env" validate:"required"`
	Server      Server   `json:"server" validate:"required"`
	Upstream    Upstream `json:"upstream" validate:"required"`
}

type Upstream struct {
	EventSerivce EventService `json:"event_service" validate:"required"`
	AuthService  AuthService  `json:"auth_service" validate:"required"`
}

type EventService struct {
	Url string `json:"url" validate:"required"`
}

type AuthService struct {
	Url string `json:"url" validate:"required"`
}

type Server struct {
	Port            string `json:"port" validate:"required"`
	WaitingShutdown int    `json:"waiting_shutdown" validate:"required,min=1"`
	HeaderTimeout   int    `json:"header_timeout" validate:"required,min=1"`
	ReadTimeout     int    `json:"read_timeout" validate:"required,min=1"`
	WriteTimeout    int    `json:"write_timeout" validate:"required,min=1"`
	IdleTimeout     int    `json:"idle_timeout" validate:"required,min=1"`
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
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("settings: validation err: %w", err)
	}

	return nil
}
