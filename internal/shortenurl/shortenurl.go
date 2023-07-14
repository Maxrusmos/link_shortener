package shortenurl

import (
	"crypto/md5"
	"encoding/hex"
)

func Shortener(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}
