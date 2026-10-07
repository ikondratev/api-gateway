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
	Environment string `json:"env"`
	Server 		Server `json:"server"`
}

type Server struct {
	Host 			string 	`json:"host"`
	Port 			string	`json:"port"`
	WaitingShutdown int 	`json:"waiting_shutdown"`
	HeaderTimeout 	int		`json:"header_timeout"`
	ReadTimeout		int		`json:"read_timeout"`
	WriteTimeout	int		`json:"write_teimeout"`
	IdleTimeout		int		`json:"idle_timeout"`

}

func New(env string) (*Settings, error){
	path := fmt.Sprintf("%s/%s_%s", settingsPath, env, settingsName)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("LoadFileErorr: %w", err)
	}

	var set *Settings
	if err := json.Unmarshal(data, &set); err != nil {
		return nil, fmt.Errorf("Unnable unmarshal config: %w", err)
	}

	return set, nil
}