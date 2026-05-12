package MyEncryption

import "fmt"

func Encryption(originalData, key []byte) (ret []byte, err error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("key is empty")
	}

	// key 仍然限制 ASCII，避免 key 本身引入复杂情况
	for _, c := range key {
		if c > 127 {
			return nil, fmt.Errorf("the key text has more than 127 characters")
		}
	}

	ret = make([]byte, 0, len(originalData)*2)

	keyIndex := 0
	for _, c := range originalData {
		if keyIndex >= len(key) {
			keyIndex = 0
		}

		v := uint16(c) + uint16(key[keyIndex])

		// 固定 2 byte，大端存储
		ret = append(ret, byte(v>>8), byte(v))

		keyIndex++
	}

	return ret, nil
}

func Decrypt(cipherText, key []byte) (ret []byte, err error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("key is empty")
	}

	if len(cipherText)%2 != 0 {
		return nil, fmt.Errorf("invalid cipher text length")
	}

	for _, c := range key {
		if c > 127 {
			return nil, fmt.Errorf("the key text has more than 127 characters")
		}
	}

	ret = make([]byte, 0, len(cipherText)/2)

	keyIndex := 0
	for i := 0; i < len(cipherText); i += 2 {
		if keyIndex >= len(key) {
			keyIndex = 0
		}

		v := uint16(cipherText[i])<<8 | uint16(cipherText[i+1])
		k := uint16(key[keyIndex])

		if v < k {
			return nil, fmt.Errorf("invalid cipher text")
		}

		ret = append(ret, byte(v-k))

		keyIndex++
	}

	return ret, nil
}
