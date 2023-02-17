package MyEncryption

import "testing"

func TestEncryption(t *testing.T) {
	key := []byte("Shehuizhu1hao")
	originalData := "H3tTAvMQk2vmhE+PCkDvocUeTyoRgHu72Bm6hYedImWfoJIkmbgUUz/OTRYBlha0qmZ5kijSnPUP86FJ/rlI7ii0JnzZ4/7ngb2bd8cg62FykaqpCyoHg2RKl6ugmKXSjcjw7l+9RXDrA21I9da6F0KRMj6e4G/L3f6cvhc4LxM="
	ret, err := Encryption([]byte(originalData), key)
	if err == nil {
		sRet := string(ret)
		t.Log(string(ret), len(ret), len(sRet))
	} else {
		t.Error(err)
	}
	original := Decrypt(ret, key)
	if string(original) != originalData {
		t.Errorf("excepted = %v, get = %v", originalData, original)
	}
}

func TestDecrypt(t *testing.T) {

}
