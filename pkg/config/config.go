package config

type Config struct {
	Cache          string          `json:"cache"`
	Port           string          `json:"port"`
	Handlers       []Handler       `json:"handlers"`
	ExternalAssets []ExternalAsset `json:"externalAssets"`
}

type Handler struct {
	Path  string `json:"path"`
	Asset string `json:"asset"`
}

type ExternalAsset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
