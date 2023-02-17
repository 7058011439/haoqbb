package AES

import (
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	origData := "Hello World" // 待加密的数据
	//key := "ABCDEFGHIJKLMNOP" // 加密的密钥
	key := "Shun7ziranShun7" // 加密的密钥
	log.Println("原文：", origData)

	log.Println("------------------ CBC模式 --------------------")
	encrypted := EncryptCBC(origData, key)
	log.Println("密文(base64)：", encrypted)
	decrypted := DecryptCBC(encrypted, key)
	log.Println("解密结果：", decrypted)

	log.Println("------------------ ECB模式 --------------------")
	encrypted = EncryptECB(origData, key)
	log.Println("密文(base64)：", encrypted)
	decrypted = DecryptECB(encrypted, key)
	log.Println("解密结果：", decrypted)

	log.Println("------------------ CFB模式 --------------------")
	encrypted = EncryptCFB(origData, key)
	log.Println("密文(base64)：", encrypted)
	decrypted = DecryptCFB(encrypted, key)
	log.Println("解密结果：", decrypted)
}
