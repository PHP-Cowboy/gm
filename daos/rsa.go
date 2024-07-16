package daos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// RSAKeys represents RSA public and private keys
type RSAKeys struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSAKeysFromStrings creates a new RSAKeys instance from public and private key strings
// (Note: In real-world usage, you would typically load keys from files or other sources)
func NewRSAKeysFromStrings(publicKeyStr, privateKeyStr string) (*RSAKeys, error) {
	publicKey, err := parsePublicKey(publicKeyStr)
	if err != nil {
		return nil, err
	}

	privateKey, err := parsePrivateKey(privateKeyStr)
	if err != nil {
		return nil, err
	}

	return &RSAKeys{publicKey: publicKey, privateKey: privateKey}, nil
}

// parsePublicKey parses a PEM-encoded public key string into an *rsa.PublicKey
func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}

	return rsaPub, nil
}

// parsePrivateKey parses a PEM-encoded private key string into an *rsa.PrivateKey
func parsePrivateKey(privateKeyStr string) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKeyStr))

	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing the private key")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		var parsePkAny any
		parsePkAny, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		rsaPrivate, ok := parsePkAny.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not of type RSA")
		}

		privateKey = rsaPrivate
	}

	return privateKey, nil
}

// Encrypt 使用私钥对数据进行RSA加密
func privateEncrypt(privateKey *rsa.PrivateKey, plaintext []byte) (signature []byte, err error) {
	// 签名
	hash := sha256.New()
	hash.Write(plaintext)
	hashed := hash.Sum(nil)

	signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	return
}
