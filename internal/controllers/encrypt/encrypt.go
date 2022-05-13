package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
)

const (
	randChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" //隨機字串字符源
)

var aesIV []byte  //固定16
var aesKey string //固定32

func init() {
	NewAesKey() //初始化AES Key
}

//產生固定大小的隨機字串
func randString(number int) string {
	cache := make([]byte, number)
	for index := range cache {
		cache[index] = randChars[rand.Intn(len(randChars))]
	}
	return string(cache)
}

func NewAesKey() {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()     //上鎖並免檔案清單存取發生崩潰
	aesKey = randString(32) //隨機aesKey
}

//加密
func Encode(inputText string) (outputCode string, err error) {

	//需要去加密的字串
	plaintext := []byte(inputText)

	//建立加密演算法 aes
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("Error: NewCipher(%d bytes) = %s", len(aesKey), err)
	}

	//加密字串
	if len(aesIV) != 16 {
		aesIV = []byte(randString(16)) //隨機IV
	}

	cfb := cipher.NewCFBEncrypter(block, aesIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return fmt.Sprintf("%x", ciphertext), nil

}

//解密
func Decode(inputCode string) (outputText string, err error) {

	byteString, err := hex.DecodeString(inputCode)
	if err != nil {
		return "", fmt.Errorf("Error code : %v = %s", inputCode, err)
	}
	ciphertext := []byte(byteString)

	if len(aesIV) != 16 {
		aesIV = []byte("1234567890123456")
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("Error: NewCipher(%d bytes) = %s", len(aesKey), err)
	}

	// 解密字串
	cfbdec := cipher.NewCFBDecrypter(block, aesIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	return string(plaintextCopy), nil
}
