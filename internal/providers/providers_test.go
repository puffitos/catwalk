package providers

import (
	"slices"
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

func TestBedrockProvider_NoRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "global.") {
			t.Errorf("expected only global models without AWS_REGION, got %q", m.ID)
		}
	}
	if !hasPrefix(p.DefaultLargeModelID, "global.") {
		t.Errorf("DefaultLargeModelID = %q, want global. prefix", p.DefaultLargeModelID)
	}
	if !hasPrefix(p.DefaultSmallModelID, "global.") {
		t.Errorf("DefaultSmallModelID = %q, want global. prefix", p.DefaultSmallModelID)
	}
}

func TestBedrockProvider_EURegion(t *testing.T) {
	t.Setenv("AWS_REGION", "eu-central-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	var ids []string
	for _, m := range p.Models {
		ids = append(ids, m.ID)
		if !hasPrefix(m.ID, "eu.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for eu-central-1 (want eu. or global.)", m.ID)
		}
	}

	if !hasPrefix(p.DefaultLargeModelID, "eu.") {
		t.Errorf("DefaultLargeModelID = %q, want eu. prefix", p.DefaultLargeModelID)
	}
	if !hasPrefix(p.DefaultSmallModelID, "eu.") {
		t.Errorf("DefaultSmallModelID = %q, want eu. prefix", p.DefaultSmallModelID)
	}
}

func TestBedrockProvider_USRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "us.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for us-east-1 (want us. or global.)", m.ID)
		}
	}
	if !hasPrefix(p.DefaultLargeModelID, "us.") {
		t.Errorf("DefaultLargeModelID = %q, want us. prefix", p.DefaultLargeModelID)
	}
}

func TestBedrockProvider_JapanRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "ap-northeast-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "jp.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for ap-northeast-1 (want jp. or global.)", m.ID)
		}
	}
}

func TestBedrockProvider_AustraliaRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "ap-southeast-2")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "au.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for ap-southeast-2 (want au. or global.)", m.ID)
		}
	}
}

func TestBedrockProvider_APACRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "ap-southeast-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "apac.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for ap-southeast-1 (want apac. or global.)", m.ID)
		}
	}
}

func TestBedrockProvider_CanadaRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "ca-central-1")
	t.Setenv("AWS_DEFAULT_REGION", "")

	p := bedrockProvider()

	for _, m := range p.Models {
		if !hasPrefix(m.ID, "us.") && !hasPrefix(m.ID, "global.") {
			t.Errorf("unexpected model %q for ca-central-1 (want us. or global.)", m.ID)
		}
	}
}

func TestBedrockRegionPrefix(t *testing.T) {
	tests := []struct {
		region string
		want   string
	}{
		{"us-east-1", "us."},
		{"us-west-2", "us."},
		{"ca-central-1", "us."},
		{"eu-central-1", "eu."},
		{"eu-west-1", "eu."},
		{"ap-northeast-1", "jp."},
		{"ap-southeast-2", "au."},
		{"ap-southeast-1", "apac."},
		{"ap-northeast-2", "apac."},
		{"ap-south-1", "apac."},
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

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
