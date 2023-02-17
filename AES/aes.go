package AES

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func formatKey(key string) string {
	l := len(key)
	switch {
	case l <= 16:
		return key + string(bytes.Repeat([]byte{0x0}, 16-len(key)))
	case l <= 24:
		return key + string(bytes.Repeat([]byte{0x0}, 24-len(key)))
	case l <= 32:
		return key + string(bytes.Repeat([]byte{0x0}, 32-len(key)))
	default:
		return key[:32]
	}
}

func EncryptCBC(origData string, key string) string {
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	key = formatKey(key)
	block, _ := aes.NewCipher([]byte(key))
	blockSize := block.BlockSize()                                      // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                        // 补全码
	blockMode := cipher.NewCBCEncrypter(block, []byte(key[:blockSize])) // 加密模式
	encrypted := make([]byte, len(origData))                            // 创建数组
	blockMode.CryptBlocks(encrypted, []byte(origData))                  // 加密
	return base64.StdEncoding.EncodeToString(encrypted)
}

func DecryptCBC(encrypted string, key string) string {
	key = formatKey(key)
	encryptedByte, _ := base64.StdEncoding.DecodeString(encrypted)
	block, _ := aes.NewCipher([]byte(key))                              // 分组秘钥
	blockSize := block.BlockSize()                                      // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, []byte(key[:blockSize])) // 加密模式
	decrypted := make([]byte, len(encryptedByte))                       // 创建数组
	blockMode.CryptBlocks(decrypted, encryptedByte)                     // 解密
	return pkcs5UnPadding(string(decrypted))                            // 去除补全码
}

func pkcs5Padding(ciphertext string, blockSize int) string {
	padding := blockSize - len(ciphertext)%blockSize
	padText := string(bytes.Repeat([]byte{byte(padding)}, padding))
	return ciphertext + padText
}

func pkcs5UnPadding(origData string) string {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func EncryptECB(origData string, key string) string {
	key = formatKey(key)
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}

func DecryptECB(encrypted string, key string) string {
	key = formatKey(key)
	encryptedByte, _ := base64.StdEncoding.DecodeString(encrypted)
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted := make([]byte, len(encryptedByte))
	for bs, be := 0, cipher.BlockSize(); bs < len(encryptedByte); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encryptedByte[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return string(decrypted[:trim])
}

func generateKey(key string) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func EncryptCFB(origData string, key string) string {
	key = formatKey(key)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	encrypted := make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], []byte(origData))
	return base64.StdEncoding.EncodeToString(encrypted)
}

func DecryptCFB(encrypted string, key string) string {
	key = formatKey(key)
	encryptedByte, _ := base64.StdEncoding.DecodeString(encrypted)
	block, _ := aes.NewCipher([]byte(key))
	if len(encryptedByte) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encryptedByte[:aes.BlockSize]
	encryptedByte = encryptedByte[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedByte, encryptedByte)
	return string(encryptedByte)
}
