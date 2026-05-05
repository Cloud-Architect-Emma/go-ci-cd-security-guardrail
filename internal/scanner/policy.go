package scanner

import (
	"encoding/json"
	"os"
)

type Rule struct {
	ID       string `json:"id"`
	Pattern  string `json:"pattern"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type Policy struct {
	Rules []Rule `json:"rules"`
}

func LoadPolicy(path string) (Policy, error) {
	var policy Policy
	data, err := os.ReadFile(path)
	if err != nil {
		return policy, err
	}

	err = json.Unmarshal(data, &policy)
	return policy, err
}