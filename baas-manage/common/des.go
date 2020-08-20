package common

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"github.com/paybf/baasmanager/baas-gateway/config"
)

//解密
func MyDESDecrypt(data string) string {
	key := config.Config.GetString("DES.key")
	crypted, _ := base64.StdEncoding.DecodeString(data)
	block, _ := des.NewCipher([]byte(key))
	blockMode := cipher.NewCBCDecrypter(block, []byte(key))
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
