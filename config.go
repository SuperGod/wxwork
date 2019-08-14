package wxwork

import (
	"encoding/json"
	"os"
	"path"
)

// Config wxwork config
type Config struct {
	CropID string `json:"cropid"`
	Agent  int    `json:"agentid"`
	Secret string `json:"secret"`
}

func GetAPI() (a *API, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configFile := path.Join(home, ".config", "wxwork.json")
	return GetAPIFromFile(configFile)
}

func GetAPIFromFile(configFile string) (a *API, err error) {
	f, err := os.Open(configFile)
	if err != nil {
		return
	}
	var cfg Config
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&cfg)
	if err != nil {
		return
	}
	a = NewAPI(cfg.CropID, cfg.Agent, cfg.Secret)
	return
}
