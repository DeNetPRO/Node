package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256String(data []byte) string {
	return hex.EncodeToString(Sha256(data))


}

func Sha256(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}


