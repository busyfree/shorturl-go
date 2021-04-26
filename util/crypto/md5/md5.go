package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func EncryptString(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptByte(text []byte) string {
	has := md5.Sum(text)
	return fmt.Sprintf("%x", has)
}
