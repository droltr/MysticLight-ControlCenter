package profile

import (
	"encoding/json"
	"fmt"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type Document struct {
	Name      string             `json:"name"`
	Ownership []domain.Ownership `json:"ownership"`
}

func ParseJSON(data []byte, knownProviders map[string]struct{}) (domain.Profile, error) {
	var document Document
	if err := json.Unmarshal(data, &document); err != nil {
		return domain.Profile{}, fmt.Errorf("parse profile JSON: %w", err)
	}
	profile := domain.Profile{
		Name:      document.Name,
		Ownership: document.Ownership,
	}
	if err := profile.Validate(knownProviders); err != nil {
		return domain.Profile{}, fmt.Errorf("validate profile: %w", err)
	}
	return profile, nil
}
