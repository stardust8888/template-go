package controller

import (
	"fmt"
	"testing"

	"template-go/utils"
)

func TestEncrypt(t *testing.T) {
	key, err := cryptKey()
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, err := utils.AESGCMEncryptHex([]byte("Qwer12356"), key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ciphertext)
}

func TestDecrypt(t *testing.T) {
	key, err := cryptKey()
	if err != nil {
		t.Fatal(err)
	}

	plaintext, err := utils.AESDecryptHex("6d7e9f4154722facc9ca50ac1678581657977cce0ed3fa05dbbaa360c27514c011d2acca5e", key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(plaintext))
}
