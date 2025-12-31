package cryptx

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	plaintext := "xxxxx"
	salt := "HWVOFkGgPTryzICwd7qnJaZR9KQ2i8xe"
	
	ciphertext := PasswordEncrypt(salt, plaintext)
	t.Logf("ciphertext: %s", ciphertext)
}
