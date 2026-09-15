package profile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

func TestSaveAndLoadFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profile.json")
	profile := domain.Profile{
		Name:      "quiet",
		Ownership: []domain.Ownership{{Resource: domain.ResourceRGB, Provider: "openrgb"}},
	}
	if err := SaveFile(path, profile); err != nil {
		t.Fatalf("failed to save profile: %v", err)
	}
	loaded, err := LoadFile(path, map[string]struct{}{"openrgb": {}})
	if err != nil {
		t.Fatalf("failed to load profile: %v", err)
	}
	if loaded.Name != profile.Name || len(loaded.Ownership) != 1 {
		t.Fatalf("unexpected loaded profile: %#v", loaded)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat profile: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("profile file permissions are too broad: %o", info.Mode().Perm())
	}
}
