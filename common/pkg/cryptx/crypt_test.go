package cryptx

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	plaintext := ""
	salt := "="

	ciphertext := PasswordEncrypt(salt, plaintext)
	t.Logf("ciphertext: %s", ciphertext)
}
