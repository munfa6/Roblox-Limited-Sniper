package data

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func (s *Config) ToJson() (Data []byte, err error) {
	return json.MarshalIndent(s, "", "  ")
}

func (config *Config) SaveConfig() {
	if Json, err := config.ToJson(); err == nil {
		WriteFile("config.json", string(Json))
	}
}

func (s *Config) LoadState() {
	data, err := ReadFile("config.json")
	if err != nil {
		s.LoadFromFile()
		s.Percentage = 10.0
		s.SaveConfig()
		return
	}

	json.Unmarshal([]byte(data), s)
	s.LoadFromFile()
}

func (c *Config) LoadFromFile() {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		jsonFile, _ = os.Create("config.json")
	}

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c)
}

func WriteFile(path string, content string) {
	os.WriteFile(path, []byte(content), 0644)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func CheckForValidFile(input string) bool {
	_, err := os.Stat(input)
	return errors.Is(err, os.ErrNotExist)
}
