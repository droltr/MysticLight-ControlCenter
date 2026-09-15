package profile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func LoadFile(path string, knownProviders map[string]struct{}) (domain.Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("read profile file: %w", err)
	}
	return ParseJSON(data, knownProviders)
}

func SaveFile(path string, profile domain.Profile) error {
	data, err := marshalJSON(profile)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write profile file: %w", err)
	}
	return nil
}

func marshalJSON(profile domain.Profile) ([]byte, error) {
	data, err := json.MarshalIndent(Document{
		Name:      profile.Name,
		Ownership: profile.Ownership,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal profile JSON: %w", err)
	}
	return append(data, '\n'), nil
}
