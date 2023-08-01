package license

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/zcalusic/sysinfo"
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

func VerifyLicense(path string) bool {
	var si sysinfo.SysInfo
	si.GetSysInfo()
	var licstr string
	if len(si.Network) != 0 {
		licstr = fmt.Sprintf("%s%d%s%sMid%s", si.Network[0].MACAddress, si.CPU.Speed, si.BIOS.Vendor, si.Board.Version, si.Node.MachineID)
	} else {
		licstr = fmt.Sprintf("%s%d%s%sMid%s", "Aa:Bb:Cc:dD:eE:fF", si.CPU.Speed, si.BIOS.Vendor, si.Board.Version, si.Node.MachineID)
	}
	data := []byte(licstr)
	md := md5.Sum(data)
	localMachineID := fmt.Sprintf("%x", md)
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	fileMachinID, err := DecryptByAes(string(buf), AESkey)
	if err != nil {
		return false
	}
	fStr := string(fileMachinID)[:len(localMachineID)]
	if localMachineID == fStr {
		return true
	}
	return false
}

func EncryptAES(txt []byte) (string, error) {
	res, err := AesEncrypt(txt, AESkey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptByAes(data string, key []byte) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, dataByte)
	crypted = pkcs7Padding(crypted, len(crypted))
	if err != nil {
		return nil, err
	}
	return crypted, nil
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
