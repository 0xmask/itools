package icrypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

type Encrypt struct {
	data []byte
}

func Sha1(data []byte) *Encrypt {
	h := sha1.Sum(data)
	return &Encrypt{data: h[:]}
}

func Sha256(data []byte) *Encrypt {
	h := sha256.Sum256(data)
	return &Encrypt{data: h[:]}
}

func Sha512(data []byte) *Encrypt {
	h := sha512.Sum512(data)
	return &Encrypt{data: h[:]}
}

func Md5(data []byte) *Encrypt {
	h := md5.Sum(data)
	return &Encrypt{data: h[:]}
}

func HmacSha256(data, key []byte) *Encrypt {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(data)
	return &Encrypt{data: h.Sum(nil)}
}

func (t *Encrypt) ToHex() string {
	return hex.EncodeToString(t.data)
}

func (t *Encrypt) ToByte() []byte {
	return t.data
}
