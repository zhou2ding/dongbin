package license

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var AESkey = []byte("JLs-kl2,+W)#(%6N*QK)r85kbt0D7$N3")

func GenerateLicense(path string) {
	srcFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	info, _ := srcFile.Stat()
	buf := make([]byte, info.Size())
	srcFile.Read(buf)
	license, _ := EncryptAES(buf)

	dstFile, err := os.Create("./license.txt")
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	dstFile.WriteString(license)
	dstFile.Sync()
}

func EncryptAES(txt []byte) (string, error) {
	res, err := AesEncrypt(txt, AESkey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
