package profile

import "testing"

func TestParseJSONValidatesProfile(t *testing.T) {
	profile, err := ParseJSON([]byte(`{
  "name": "quiet",
  "ownership": [{"resource": "rgb", "provider": "openrgb"}]
}`), map[string]struct{}{"openrgb": {}})
	if err != nil {
		t.Fatalf("failed to parse profile: %v", err)
	}
	if profile.Name != "quiet" {
		t.Fatalf("unexpected profile name: %q", profile.Name)
	}
}

func TestParseJSONRejectsInvalidProfile(t *testing.T) {
	_, err := ParseJSON([]byte(`{"name":"unsafe","ownership":[{"Resource":"rgb","Provider":"missing"}]}`), map[string]struct{}{"openrgb": {}})
	if err == nil {
		t.Fatal("expected invalid provider to be rejected")
	}
}
