package providers

import (
	"slices"
	"strings"
	"testing"
)

func TestValidDefaultModels(t *testing.T) {
	for _, p := range GetAll() {
		t.Run(p.Name, func(t *testing.T) {
			var modelIds []string
			for _, m := range p.Models {
				modelIds = append(modelIds, m.ID)
			}
			if !slices.Contains(modelIds, p.DefaultLargeModelID) {
				t.Errorf("Default large model %q not found in provider %q", p.DefaultLargeModelID, p.Name)
			}
			if !slices.Contains(modelIds, p.DefaultSmallModelID) {
				t.Errorf("Default small model %q not found in provider %q", p.DefaultSmallModelID, p.Name)
			}
		})
	}
}

func TestBedrockProvider(t *testing.T) {
	totalModels := len(loadProviderFromConfig(bedrockConfig).Models)

	tests := []struct {
		name           string
		region         string
		wantPrefix     string
		wantDefaultPfx string
		wantModels     int
	}{
		{"no region falls back to global", "", "global.", "global.", totalModels - 2},
		{"unknown region falls back to global", "sa-east-1", "global.", "global.", totalModels - 2},
		{"eu-central-1", "eu-central-1", "eu.", "eu.", totalModels - 2},
		{"us-east-1", "us-east-1", "us.", "us.", totalModels},
		{"ca-central-1 maps to us", "ca-central-1", "us.", "us.", totalModels},
		{"ap-northeast-1 maps to jp", "ap-northeast-1", "jp.", "jp.", totalModels - 2},
		{"ap-southeast-2 maps to au", "ap-southeast-2", "au.", "au.", totalModels - 2},
		{"ap-southeast-1 maps to apac", "ap-southeast-1", "apac.", "global.", totalModels - 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("AWS_REGION", tt.region)
			t.Setenv("AWS_DEFAULT_REGION", "")

			p := bedrockProvider()

			if len(p.Models) != tt.wantModels {
				t.Errorf("got %d models, want %d", len(p.Models), tt.wantModels)
			}

			// All model IDs must carry the expected prefix or global.
			for _, m := range p.Models {
				if !strings.HasPrefix(m.ID, tt.wantPrefix) && !strings.HasPrefix(m.ID, "global.") {
					t.Errorf("model %q has unexpected prefix for region %q", m.ID, tt.region)
				}
			}

			// Default model IDs must use the expected prefix.
			if !strings.HasPrefix(p.DefaultLargeModelID, tt.wantDefaultPfx) {
				t.Errorf("DefaultLargeModelID = %q, want %q prefix", p.DefaultLargeModelID, tt.wantDefaultPfx)
			}
			if !strings.HasPrefix(p.DefaultSmallModelID, tt.wantDefaultPfx) {
				t.Errorf("DefaultSmallModelID = %q, want %q prefix", p.DefaultSmallModelID, tt.wantDefaultPfx)
			}

			// Default model IDs must exist in the model list.
			var ids []string
			for _, m := range p.Models {
				ids = append(ids, m.ID)
			}
			if !slices.Contains(ids, p.DefaultLargeModelID) {
				t.Errorf("DefaultLargeModelID %q not found in model list", p.DefaultLargeModelID)
			}
			if !slices.Contains(ids, p.DefaultSmallModelID) {
				t.Errorf("DefaultSmallModelID %q not found in model list", p.DefaultSmallModelID)
			}
		})
	}
}

// TestBedrockConfigRegions asserts that regions for the inference
// profile mapping are configured in the bedrock configuration file
func TestBedrockConfigRegions(t *testing.T) {
	t.Parallel()

	p := loadProviderFromConfig(bedrockConfig)
	for _, m := range p.Models {
		if len(m.Regions) == 0 {
			t.Errorf("model %q has no regions configured, at least one must be defined.", m.ID)
		}
	}
}

func TestBedrockRegionPrefix(t *testing.T) {
	tests := []struct {
		region string
		want   string
	}{
		{"us-east-1", "us"},
		{"us-west-2", "us"},
		{"ca-central-1", "us"},
		{"eu-central-1", "eu"},
		{"eu-west-1", "eu"},
		{"ap-northeast-1", "jp"},
		{"ap-southeast-2", "au"},
		{"ap-southeast-1", "apac"},
		{"ap-northeast-2", "apac"},
		{"ap-south-1", "apac"},
		{"sa-east-1", ""},
		{"", ""},
		{"unknown-region", ""},
	}

	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			got := bedrockRegionPrefix(tt.region)
			if got != tt.want {
				t.Errorf("bedrockRegionPrefix(%q) = %q, want %q", tt.region, got, tt.want)
			}
		})
	}
}
