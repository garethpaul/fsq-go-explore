// Foursquare GoLang SDK
package fsq

import (
  "log"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
)

type s struct {
	Int       int
	String    string
	ByteSlice []byte
}

// Return a key based on a given VenueSearchRequest with AES encryption
func GetSearchKey(q *VenueSearchRequest) (s string) {
	a := q
	out, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	// Create the encryption algorithm
	key_text := "astaxie12798akljzmknm.ahkjkljl;k"
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		log.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
	}

	// Encrypt the key
	var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(out))
	cfb.XORKeyStream(ciphertext, out)

	key_out := base64.StdEncoding.EncodeToString([]byte(ciphertext))
	return key_out
}
