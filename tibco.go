package main

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// fixedKey is TIBCO's hardcoded 24-byte 3DES (DESede) key used by the
// ObfuscationEngine "fixed key" mode — the secrets prefixed with "#!".
// It ships inside TIBCO's own libraries, so these values are reversible by
// design. Use only for secrets from systems you administer.
var fixedKey = []byte{
	28, 167, 155, 145, 91, 143, 26, 186,
	98, 176, 233, 203, 138, 93, 173, 239,
	28, 167, 155, 145, 91, 143, 26, 186,
}

const prefix = "#!"

// Deobfuscate reverses a TIBCO fixed-key secret. "#!" prefix is optional.
// Layout after base64-decode: IV(8) || ciphertext; plaintext is UTF-16LE.
func Deobfuscate(s string) (string, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, prefix)

	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	block, err := des.NewTripleDESCipher(fixedKey)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	if len(raw) < bs {
		return "", errors.New("input too short to contain an IV")
	}
	iv, ct := raw[:bs], raw[bs:]
	if len(ct) == 0 || len(ct)%bs != 0 {
		return "", errors.New("ciphertext length is not a block multiple")
	}

	pt := make([]byte, len(ct))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(pt, ct)

	pt, err = pkcs5Unpad(pt, bs)
	if err != nil {
		return "", err
	}
	return utf16leToString(pt), nil
}

// Obfuscate produces a "#!"-prefixed TIBCO fixed-key secret with a random IV.
func Obfuscate(plain string) (string, error) {
	block, err := des.NewTripleDESCipher(fixedKey)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()

	iv := make([]byte, bs)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	pt := pkcs5Pad(stringToUTF16le(plain), bs)
	ct := make([]byte, len(pt))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ct, pt)

	out := append(append([]byte{}, iv...), ct...)
	return prefix + base64.StdEncoding.EncodeToString(out), nil
}

func pkcs5Pad(b []byte, bs int) []byte {
	n := bs - len(b)%bs
	return append(b, strings.Repeat(string(rune(n)), n)...)
}

func pkcs5Unpad(b []byte, bs int) ([]byte, error) {
	if len(b) == 0 || len(b)%bs != 0 {
		return nil, errors.New("invalid padded length")
	}
	n := int(b[len(b)-1])
	if n == 0 || n > bs || n > len(b) {
		return nil, errors.New("invalid padding")
	}
	for _, c := range b[len(b)-n:] {
		if int(c) != n {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return b[:len(b)-n], nil
}

func utf16leToString(b []byte) string {
	var sb strings.Builder
	for i := 0; i+1 < len(b); i += 2 {
		sb.WriteRune(rune(uint16(b[i]) | uint16(b[i+1])<<8))
	}
	return sb.String()
}

func stringToUTF16le(s string) []byte {
	var out []byte
	for _, r := range s {
		out = append(out, byte(r), byte(r>>8))
	}
	return out
}
