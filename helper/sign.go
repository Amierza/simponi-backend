package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateShopeeSign(partnerID, path string, timestamp int64, partnerKey string) string {
	base := fmt.Sprintf("%s%s%d", partnerID, path, timestamp)

	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(base))

	return hex.EncodeToString(h.Sum(nil))
}

func GenerateShopeeShopSign(partnerID string, path string, timestamp int64, accessToken string, shopID string, partnerKey string) string {
	base := fmt.Sprintf(
		"%s%s%d%s%s",
		partnerID,
		path,
		timestamp,
		accessToken,
		shopID,
	)

	h := hmac.New(sha256.New, []byte(partnerKey))
	h.Write([]byte(base))

	return hex.EncodeToString(h.Sum(nil))
}
