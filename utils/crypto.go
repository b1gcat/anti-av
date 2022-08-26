package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func Kek(src []byte) []byte {
	h := sha256.New()
	h.Write(src)
	return h.Sum(nil)[:4]
}

func Crypt(key, src []byte) ([]byte, error) {
	k := hex.EncodeToString(key[:8])

	iv := make([]byte, 8)
	rand.Read(iv)
	i := hex.EncodeToString(iv[:8])

	x, err := aes.NewCipher([]byte(k))
	if err != nil {
		return nil, err
	}
	src = PKCS7Padding(src, x.BlockSize())
	dst := make([]byte, len(src)+len(key)+len(iv))

	copy(dst, key)
	copy(dst[len(key):], iv)

	mode := cipher.NewCBCEncrypter(x, []byte(i))
	mode.CryptBlocks(dst[len(key)+len(iv):], src)
	//key+iv+e_data
	return dst, nil
}

func DeCrypt(src []byte) ([]byte, error) {
	key := hex.EncodeToString(src[:8])
	iv := hex.EncodeToString(src[8:16])
	src = src[16:]
	x, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(x, []byte(iv))
	dst := make([]byte, len(src))
	mode.CryptBlocks(dst, src)
	dst = PKCS7UnPadding(dst)
	return dst, nil
}

//PKCS7Padding say ...
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

//PKCS7UnPadding 使用PKCS7进行填充 复原
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return origData
	}
	unPadding := int(origData[length-1])
	if unPadding > length {
		return origData
	}
	return origData[:(length - unPadding)]
}
