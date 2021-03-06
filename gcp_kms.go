package kiya

import (
	"encoding/base64"
	"fmt"

	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
)

// GetDecryptedValue decrypts an encrypted value via Google KMS
func GetDecryptedValue(kmsService *cloudkms.Service, target Profile, cipherText string) (string, error) {
	decryptReq := &cloudkms.DecryptRequest{
		Ciphertext: cipherText,
	}
	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		target.ProjectID,
		target.Location,
		target.Keyring,
		target.CryptoKey)
	resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(path, decryptReq).Do()
	if err != nil {
		return "", tre.New(err, "failed to decrypt", "path", path)
	}
	data, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return "", tre.New(err, "failed to base64 decode")
	}
	return string(data), nil
}

// GetEncryptedValue converts a plain text to a Google KMS encrypted text
func GetEncryptedValue(kmsService *cloudkms.Service, target Profile, plainText string) (string, error) {
	encryptReq := &cloudkms.EncryptRequest{
		Plaintext: base64.RawURLEncoding.EncodeToString([]byte(plainText)),
	}
	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		target.ProjectID,
		target.Location,
		target.Keyring,
		target.CryptoKey)
	resp, err := kmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(path, encryptReq).Do()
	if err != nil {
		return "", tre.New(err, "failed to encrypt")
	}
	return resp.Ciphertext, nil
}
