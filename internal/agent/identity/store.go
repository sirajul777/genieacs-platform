package identity

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Credentials struct { AgentID string `json:"agent_id"`; Token string `json:"token"` }

func Load(path string) (Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil { return Credentials{}, err }
	var c Credentials
	if err := json.Unmarshal(data, &c); err != nil { return Credentials{}, err }
	if c.AgentID == "" || c.Token == "" { return Credentials{}, errors.New("invalid agent credentials") }
	return c, nil
}

func Save(path string, c Credentials) error {
	if c.AgentID == "" || c.Token == "" { return errors.New("agent credentials are required") }
	if dir := filepath.Dir(path); dir != "." { if err := os.MkdirAll(dir, 0700); err != nil { return err } }
	data, err := json.Marshal(c)
	if err != nil { return err }
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil { return err }
	if err := os.Rename(tmp, path); err != nil { _ = os.Remove(tmp); return err }
	return nil
}
