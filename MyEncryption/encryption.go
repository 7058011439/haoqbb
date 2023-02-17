package MyEncryption

import "fmt"

func Encryption(originalData, key []byte) (ret []byte, err error) {
	for _, c := range originalData {
		if c > 127 {
			return nil, fmt.Errorf("the original text has more than 127 characters")
		}
	}

	for _, c := range key {
		if c > 127 {
			return nil, fmt.Errorf("the key text has more than 127 characters")
		}
	}

	keyIndex := 0
	for _, c := range originalData {
		if keyIndex >= len(key) {
			keyIndex = 0
		}
		ret = append(ret, c+key[keyIndex])
	}

	return
}

func Decrypt(cipherText, key []byte) (ret []byte) {
	keyIndex := 0
	for _, c := range cipherText {
		if keyIndex >= len(key) {
			keyIndex = 0
		}
		ret = append(ret, c-key[keyIndex])
	}

	return
}
