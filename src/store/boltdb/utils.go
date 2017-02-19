package boltdb

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"log"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

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

func decodeToStringInterface(b []byte, d interface{}) error {
	dec := msgpack.NewDecoder(bytes.NewBuffer(b))
	dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{}, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}

			mv, err := d.DecodeInterface()
			if err != nil {
				return nil, err
			}

			m[mk] = mv
		}
		return m, nil
	}
	return dec.Decode(&d)
}
