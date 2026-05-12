package platforms

import "os"

type ShopeeConfig struct {
	PartnerID   string
	PartnerKey  string
	BaseURL     string
	RedirectURL string
}

func NewShopeeConfig() ShopeeConfig {
	return ShopeeConfig{
		PartnerID:   os.Getenv("SHOPEE_PARTNER_ID"),
		PartnerKey:  os.Getenv("SHOPEE_PARTNER_KEY"),
		BaseURL:     os.Getenv("SHOPEE_BASE_URL"),
		RedirectURL: os.Getenv("SHOPEE_REDIRECT_URL"),
	}
}

func (sc *ShopeeConfig) GetShopeeConfig() *ShopeeConfig {
	return sc
}
