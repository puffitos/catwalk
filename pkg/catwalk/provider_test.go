package catwalk

import (
	"testing"
)

func TestProvider_PrefixModelIDs(t *testing.T) {
	t.Parallel()

	t.Run("prefixes all IDs", func(t *testing.T) {
		t.Parallel()
		p := Provider{
			DefaultLargeModelID: "anthropic.claude-sonnet-4-6",
			DefaultSmallModelID: "anthropic.claude-haiku-4-5-20251001-v1:0",
			Models: []Model{
				{ID: "anthropic.claude-sonnet-4-6"},
				{ID: "anthropic.claude-haiku-4-5-20251001-v1:0"},
			},
		}

		p.PrefixModelIDs("eu.")

		if p.DefaultLargeModelID != "eu.anthropic.claude-sonnet-4-6" {
			t.Errorf("DefaultLargeModelID = %q, want %q", p.DefaultLargeModelID, "eu.anthropic.claude-sonnet-4-6")
		}
		if p.DefaultSmallModelID != "eu.anthropic.claude-haiku-4-5-20251001-v1:0" {
			t.Errorf("DefaultSmallModelID = %q, want %q", p.DefaultSmallModelID, "eu.anthropic.claude-haiku-4-5-20251001-v1:0")
		}
		if p.Models[0].ID != "eu.anthropic.claude-sonnet-4-6" {
			t.Errorf("Models[0].ID = %q, want %q", p.Models[0].ID, "eu.anthropic.claude-sonnet-4-6")
		}
		if p.Models[1].ID != "eu.anthropic.claude-haiku-4-5-20251001-v1:0" {
			t.Errorf("Models[1].ID = %q, want %q", p.Models[1].ID, "eu.anthropic.claude-haiku-4-5-20251001-v1:0")
		}
	})

	t.Run("empty prefix is a no-op", func(t *testing.T) {
		t.Parallel()
		p := Provider{
			DefaultLargeModelID: "anthropic.claude-sonnet-4-6",
			Models:              []Model{{ID: "anthropic.claude-sonnet-4-6"}},
		}

		p.PrefixModelIDs("")

		if p.DefaultLargeModelID != "anthropic.claude-sonnet-4-6" {
			t.Errorf("DefaultLargeModelID changed to %q", p.DefaultLargeModelID)
		}
		if p.Models[0].ID != "anthropic.claude-sonnet-4-6" {
			t.Errorf("Models[0].ID changed to %q", p.Models[0].ID)
		}
	})

	t.Run("empty default IDs are not prefixed", func(t *testing.T) {
		t.Parallel()
		p := Provider{
			Models: []Model{{ID: "some-model"}},
		}

		p.PrefixModelIDs("us.")

		if p.DefaultLargeModelID != "" {
			t.Errorf("DefaultLargeModelID = %q, want empty", p.DefaultLargeModelID)
		}
		if p.DefaultSmallModelID != "" {
			t.Errorf("DefaultSmallModelID = %q, want empty", p.DefaultSmallModelID)
		}
		if p.Models[0].ID != "us.some-model" {
			t.Errorf("Models[0].ID = %q, want %q", p.Models[0].ID, "us.some-model")
		}
	})
}
