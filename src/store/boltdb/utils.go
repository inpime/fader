package boltdb

import (
	"crypto/sha1"
	"encoding/hex"
	"log"

	uuid "github.com/satori/go.uuid"
)

func SHA1(v interface{}) []byte {
	hasher := sha1.New()
	switch _v := v.(type) {
	case string:
		hasher.Write([]byte(_v))
	case []byte:
		hasher.Write(_v)
	default:
		log.Println("[WARNING] SHA1: not supported type (only text OR []byte)")
	}

	return hasher.Sum(nil)
}

func SHA1String(v interface{}) string {
	return hex.EncodeToString(SHA1(v))
}

// TODO: переименовать в hashFromUUID

func hashFromFile(bucketID uuid.UUID, fileName string) []byte {
	return SHA1(bucketID.String() + fileName)
}
