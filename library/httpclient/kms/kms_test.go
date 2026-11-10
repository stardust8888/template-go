package kms

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_Encrypt(t *testing.T) {
	InitClient()
	ciphertext, err := client.Encrypt(context.Background(), "test")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ciphertext)
}

func TestClient_Decrypt(t *testing.T) {
	InitClient()
	plaintext, err := client.Decrypt(context.Background(), "8c8f45e8806a4674ce1e50ff163314476c8e89f7a31eacaa9bd46f182840c20c")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(plaintext)
}
