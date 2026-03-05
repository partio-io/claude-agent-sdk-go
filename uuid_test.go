package claude

import (
	"regexp"
	"testing"
)

func TestNewUUID(t *testing.T) {
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	t.Run("format", func(t *testing.T) {
		id := newUUID()
		if !re.MatchString(id) {
			t.Errorf("UUID %q does not match v4 format", id)
		}
	})

	t.Run("uniqueness", func(t *testing.T) {
		seen := make(map[string]bool)
		for range 1000 {
			id := newUUID()
			if seen[id] {
				t.Fatalf("duplicate UUID: %s", id)
			}
			seen[id] = true
		}
	})
}
