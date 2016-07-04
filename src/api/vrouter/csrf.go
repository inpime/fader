package vrouter

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func generateCSRFToken(secret, salt []byte) string {
	h := hmac.New(sha1.New, secret)
	h.Write(salt)
	return fmt.Sprintf("%s:%s", hex.EncodeToString(h.Sum(nil)), hex.EncodeToString(salt))
}

func validateCSRFToken(token string, secret []byte) (bool, error) {
	sep := strings.Index(token, ":")
	if sep < 0 {
		return false, nil
	}
	salt, err := hex.DecodeString(token[sep+1:])
	if err != nil {
		return false, err
	}
	return token == generateCSRFToken(secret, salt), nil
}

func generateSalt(len uint8) (salt []byte, err error) {
	salt = make([]byte, len)
	_, err = rand.Read(salt)
	return
}
