package platforms

import "os"

type MLConfig struct {
	BaseURL string
}

func NewMLConfig() MLConfig {
	return MLConfig{
		BaseURL: os.Getenv("ML_API_BASE_URL"),
	}
}

func (mc *MLConfig) GetMLConfig() *MLConfig {
	return mc
}
