package catwalk

import "strings"

// Type represents the type of AI provider.
type Type string

// All the supported AI provider types.
const (
	TypeOpenAI       Type = "openai"
	TypeOpenAICompat Type = "openai-compat"
	TypeOpenRouter   Type = "openrouter"
	TypeVercel       Type = "vercel"
	TypeAnthropic    Type = "anthropic"
	TypeGoogle       Type = "google"
	TypeAzure        Type = "azure"
	TypeBedrock      Type = "bedrock"
	TypeVertexAI     Type = "google-vertex"
)

// InferenceProvider represents the inference provider identifier.
type InferenceProvider string

// All the inference providers supported by the system.
const (
	InferenceProviderOpenAI       InferenceProvider = "openai"
	InferenceProviderAnthropic    InferenceProvider = "anthropic"
	InferenceProviderSynthetic    InferenceProvider = "synthetic"
	InferenceProviderGemini       InferenceProvider = "gemini"
	InferenceProviderAzure        InferenceProvider = "azure"
	InferenceProviderBedrock      InferenceProvider = "bedrock"
	InferenceProviderVertexAI     InferenceProvider = "vertexai"
	InferenceProviderXAI          InferenceProvider = "xai"
	InferenceProviderZAI          InferenceProvider = "zai"
	InferenceProviderZhipu        InferenceProvider = "zhipu"
	InferenceProviderZhipuCoding  InferenceProvider = "zhipu-coding"
	InferenceProviderGROQ         InferenceProvider = "groq"
	InferenceProviderOpenRouter   InferenceProvider = "openrouter"
	InferenceProviderCerebras     InferenceProvider = "cerebras"
	InferenceProviderVenice       InferenceProvider = "venice"
	InferenceProviderChutes       InferenceProvider = "chutes"
	InferenceProviderHuggingFace  InferenceProvider = "huggingface"
	InferenceAIHubMix             InferenceProvider = "aihubmix"
	InferenceKimiCoding           InferenceProvider = "kimi-coding"
	InferenceProviderCopilot      InferenceProvider = "copilot"
	InferenceProviderVercel       InferenceProvider = "vercel"
	InferenceProviderMiniMax      InferenceProvider = "minimax"
	InferenceProviderMiniMaxChina InferenceProvider = "minimax-china"
	InferenceProviderIoNet        InferenceProvider = "ionet"
	InferenceProviderQiniuCloud   InferenceProvider = "qiniucloud"
	InferenceProviderAvian        InferenceProvider = "avian"
)

// Provider represents an AI provider configuration.
type Provider struct {
	Name                string            `json:"name"`
	ID                  InferenceProvider `json:"id"`
	APIKey              string            `json:"api_key,omitempty"`
	APIEndpoint         string            `json:"api_endpoint,omitempty"`
	Type                Type              `json:"type,omitempty"`
	DefaultLargeModelID string            `json:"default_large_model_id,omitempty"`
	DefaultSmallModelID string            `json:"default_small_model_id,omitempty"`
	Models              []Model           `json:"models,omitempty"`
	DefaultHeaders      map[string]string `json:"default_headers,omitempty"`
}

// ModelOptions stores extra options for models.
type ModelOptions struct {
	Temperature      *float64       `json:"temperature,omitempty"`
	TopP             *float64       `json:"top_p,omitempty"`
	TopK             *int64         `json:"top_k,omitempty"`
	FrequencyPenalty *float64       `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64       `json:"presence_penalty,omitempty"`
	ProviderOptions  map[string]any `json:"provider_options,omitempty"`
}

// Model represents an AI model configuration.
type Model struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`
	CostPer1MIn            float64      `json:"cost_per_1m_in"`
	CostPer1MOut           float64      `json:"cost_per_1m_out"`
	CostPer1MInCached      float64      `json:"cost_per_1m_in_cached"`
	CostPer1MOutCached     float64      `json:"cost_per_1m_out_cached"`
	ContextWindow          int64        `json:"context_window"`
	DefaultMaxTokens       int64        `json:"default_max_tokens"`
	CanReason              bool         `json:"can_reason"`
	ReasoningLevels        []string     `json:"reasoning_levels,omitempty"`
	DefaultReasoningEffort string       `json:"default_reasoning_effort,omitempty"`
	SupportsImages         bool         `json:"supports_attachments"`
	Options                ModelOptions `json:"options"`
	// Regions lists the inference profile prefixes available for this model
	// (e.g. "us", "eu", "global"). Used by providers that require
	// region-specific model identifiers, such as AWS Bedrock cross-region
	// inference profiles. The actual inference profile ID is derived by
	// prepending the prefix and a "." separator to the model ID.
	Regions []string `json:"regions,omitempty"`
}

// PrefixModelIDs prepends prefix to every model ID in the provider as well
// as to DefaultLargeModelID and DefaultSmallModelID. This keeps the default
// model references in sync when a consumer (e.g. crush) adds a cross-region
// inference profile prefix for Bedrock ("us.", "eu.", "ap.").
func (p *Provider) PrefixModelIDs(prefix string) {
	if prefix == "" {
		return
	}
	for i := range p.Models {
		p.Models[i].ID = prefix + p.Models[i].ID
	}
	if p.DefaultLargeModelID != "" {
		p.DefaultLargeModelID = prefix + p.DefaultLargeModelID
	}
	if p.DefaultSmallModelID != "" {
		p.DefaultSmallModelID = prefix + p.DefaultSmallModelID
	}
}

// ApplyBedrockRegion applies the AWS cross-region inference profile prefix
// for the given region to all model IDs and default model references. For
// example, a region of "eu-central-1" prepends "eu." so that
// "anthropic.claude-sonnet-4-6" becomes "eu.anthropic.claude-sonnet-4-6".
// Regions that do not map to a known prefix are left unchanged.
func (p *Provider) ApplyBedrockRegion(region string) {
	p.PrefixModelIDs(bedrockRegionPrefix(region))
}

// bedrockRegionPrefix returns the cross-region inference profile prefix for a
// given AWS region (e.g. "eu-central-1" → "eu."), or an empty string when the
// region does not map to a known prefix.
func bedrockRegionPrefix(region string) string {
	switch {
	case strings.HasPrefix(region, "us-"):
		return "us."
	case strings.HasPrefix(region, "eu-"):
		return "eu."
	case strings.HasPrefix(region, "ap-"):
		return "ap."
	default:
		return ""
	}
}

// KnownProviders returns all the known inference providers.
func KnownProviders() []InferenceProvider {
	return []InferenceProvider{
		InferenceProviderOpenAI,
		InferenceProviderSynthetic,
		InferenceProviderAnthropic,
		InferenceProviderGemini,
		InferenceProviderAzure,
		InferenceProviderBedrock,
		InferenceProviderVertexAI,
		InferenceProviderXAI,
		InferenceProviderZAI,
		InferenceProviderZhipu,
		InferenceProviderZhipuCoding,
		InferenceProviderGROQ,
		InferenceProviderOpenRouter,
		InferenceProviderCerebras,
		InferenceProviderVenice,
		InferenceProviderChutes,
		InferenceProviderHuggingFace,
		InferenceAIHubMix,
		InferenceKimiCoding,
		InferenceProviderCopilot,
		InferenceProviderVercel,
		InferenceProviderMiniMax,
		InferenceProviderMiniMaxChina,
		InferenceProviderQiniuCloud,
		InferenceProviderAvian,
	}
}

// KnownProviderTypes returns all the known inference providers types.
func KnownProviderTypes() []Type {
	return []Type{
		TypeOpenAI,
		TypeOpenAICompat,
		TypeOpenRouter,
		TypeVercel,
		TypeAnthropic,
		TypeGoogle,
		TypeAzure,
		TypeBedrock,
		TypeVertexAI,
	}
}
