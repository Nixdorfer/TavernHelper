package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/curve25519"
)

type EncryptionType string

const (
	EncryptNone       EncryptionType = ""
	EncryptSymmetric  EncryptionType = "symmetric"
	EncryptAsymmetric EncryptionType = "asymmetric"
)

type EncryptedPackage struct {
	Version   int            `json:"v"`
	Type      EncryptionType `json:"t"`
	ExpireAt  int64          `json:"e,omitempty"`
	Nonce     string         `json:"n,omitempty"`
	PublicKey string         `json:"pk,omitempty"`
	Salt      string         `json:"s,omitempty"`
	Data      string         `json:"d"`
	Signature string         `json:"sig,omitempty"`
	RootSalt  string         `json:"rs,omitempty"`
	RootNonce string         `json:"rn,omitempty"`
	RootData  string         `json:"rd,omitempty"`
}

type EncryptionConfig struct {
	Type          EncryptionType `json:"type"`
	Password      string         `json:"password,omitempty"`
	PublicKeyPEM  string         `json:"publicKey,omitempty"`
	PrivateKeyPEM string         `json:"privateKey,omitempty"`
	ExpireAt      *time.Time     `json:"expireAt,omitempty"`
}

func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(data); err != nil {
		return nil, err
	}
	if err := gzWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompressData(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	return io.ReadAll(gzReader)
}

func deriveKeyFromPassword(password string, salt []byte) []byte {
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	return h.Sum(nil)
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptAESGCM(plaintext, key []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func decryptAESGCM(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func ed25519PublicToX25519(edPub ed25519.PublicKey) ([]byte, error) {
	if len(edPub) != ed25519.PublicKeySize {
		return nil, errors.New("invalid Ed25519 public key size")
	}
	return edPub[:32], nil
}

func ed25519PrivateToX25519(edPriv ed25519.PrivateKey) ([]byte, error) {
	if len(edPriv) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid Ed25519 private key size")
	}
	h := sha256.Sum256(edPriv.Seed())
	return h[:], nil
}

func EncryptProject(data []byte, config EncryptionConfig) (string, error) {
	compressed, err := compressData(data)
	if err != nil {
		return "", fmt.Errorf("compression failed: %w", err)
	}

	pkg := EncryptedPackage{
		Version: 1,
		Type:    config.Type,
	}

	if config.ExpireAt != nil {
		pkg.ExpireAt = config.ExpireAt.Unix()
	}

	switch config.Type {
	case EncryptSymmetric:
		salt, err := generateSalt()
		if err != nil {
			return "", err
		}
		key := deriveKeyFromPassword(config.Password, salt)
		ciphertext, nonce, err := encryptAESGCM(compressed, key)
		if err != nil {
			return "", err
		}
		pkg.Salt = base64.RawURLEncoding.EncodeToString(salt)
		pkg.Nonce = base64.RawURLEncoding.EncodeToString(nonce)
		pkg.Data = base64.RawURLEncoding.EncodeToString(ciphertext)

	case EncryptAsymmetric:
		recipientPubKey, err := parseEd25519PublicKey(config.PublicKeyPEM)
		if err != nil {
			return "", fmt.Errorf("invalid recipient public key: %w", err)
		}

		recipientX25519, err := ed25519PublicToX25519(recipientPubKey)
		if err != nil {
			return "", err
		}

		var ephemeralPrivate, ephemeralPublic [32]byte
		if _, err := rand.Read(ephemeralPrivate[:]); err != nil {
			return "", err
		}
		curve25519.ScalarBaseMult(&ephemeralPublic, &ephemeralPrivate)

		var recipientX25519Array [32]byte
		copy(recipientX25519Array[:], recipientX25519)
		var sharedSecret [32]byte
		curve25519.ScalarMult(&sharedSecret, &ephemeralPrivate, &recipientX25519Array)

		h := sha256.Sum256(sharedSecret[:])
		encryptionKey := h[:]

		ciphertext, nonce, err := encryptAESGCM(compressed, encryptionKey)
		if err != nil {
			return "", err
		}

		pkg.PublicKey = base64.RawURLEncoding.EncodeToString(ephemeralPublic[:])
		pkg.Nonce = base64.RawURLEncoding.EncodeToString(nonce)
		pkg.Data = base64.RawURLEncoding.EncodeToString(ciphertext)

	default:
		pkg.Data = base64.RawURLEncoding.EncodeToString(compressed)
	}

	result, err := json.Marshal(pkg)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(result), nil
}

func DecryptProject(encryptedStr string, password string, privateKeyPEM string) ([]byte, error) {
	jsonData, err := base64.RawURLEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, fmt.Errorf("invalid encrypted data format: %w", err)
	}

	var pkg EncryptedPackage
	if err := json.Unmarshal(jsonData, &pkg); err != nil {
		return nil, fmt.Errorf("invalid package format: %w", err)
	}

	if pkg.ExpireAt > 0 && time.Now().Unix() > pkg.ExpireAt {
		return nil, errors.New("encrypted data has expired")
	}

	var compressed []byte

	switch pkg.Type {
	case EncryptSymmetric:
		if password == "" {
			return nil, errors.New("password required for symmetric decryption")
		}
		salt, err := base64.RawURLEncoding.DecodeString(pkg.Salt)
		if err != nil {
			return nil, err
		}
		nonce, err := base64.RawURLEncoding.DecodeString(pkg.Nonce)
		if err != nil {
			return nil, err
		}
		ciphertext, err := base64.RawURLEncoding.DecodeString(pkg.Data)
		if err != nil {
			return nil, err
		}
		key := deriveKeyFromPassword(password, salt)
		compressed, err = decryptAESGCM(ciphertext, key, nonce)
		if err != nil {
			return nil, errors.New("decryption failed: invalid password or corrupted data")
		}

	case EncryptAsymmetric:
		if privateKeyPEM == "" {
			return nil, errors.New("private key required for asymmetric decryption")
		}
		privateKey, err := parseEd25519PrivateKey(privateKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}

		x25519Private, err := ed25519PrivateToX25519(privateKey)
		if err != nil {
			return nil, err
		}

		ephemeralPubBytes, err := base64.RawURLEncoding.DecodeString(pkg.PublicKey)
		if err != nil {
			return nil, err
		}

		var x25519PrivateArray, ephemeralPub [32]byte
		copy(x25519PrivateArray[:], x25519Private)
		copy(ephemeralPub[:], ephemeralPubBytes)
		var sharedSecret [32]byte
		curve25519.ScalarMult(&sharedSecret, &x25519PrivateArray, &ephemeralPub)

		h := sha256.Sum256(sharedSecret[:])
		decryptionKey := h[:]

		nonce, err := base64.RawURLEncoding.DecodeString(pkg.Nonce)
		if err != nil {
			return nil, err
		}
		ciphertext, err := base64.RawURLEncoding.DecodeString(pkg.Data)
		if err != nil {
			return nil, err
		}
		compressed, err = decryptAESGCM(ciphertext, decryptionKey, nonce)
		if err != nil {
			return nil, errors.New("decryption failed: invalid key or corrupted data")
		}

	default:
		var err error
		compressed, err = base64.RawURLEncoding.DecodeString(pkg.Data)
		if err != nil {
			return nil, err
		}
	}

	return decompressData(compressed)
}

func parseEd25519PublicKey(keyStr string) (ed25519.PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(keyStr)
	if err == nil && len(decoded) == ed25519.PublicKeySize {
		return ed25519.PublicKey(decoded), nil
	}

	decoded, err = base64.RawURLEncoding.DecodeString(keyStr)
	if err == nil && len(decoded) == ed25519.PublicKeySize {
		return ed25519.PublicKey(decoded), nil
	}

	return nil, errors.New("unsupported public key format")
}

func parseEd25519PrivateKey(keyStr string) (ed25519.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(keyStr)
	if err == nil && len(decoded) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(decoded), nil
	}

	decoded, err = base64.RawURLEncoding.DecodeString(keyStr)
	if err == nil && len(decoded) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(decoded), nil
	}

	return nil, errors.New("unsupported private key format")
}

func GenerateKeyPair() (publicKey, privateKey string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	publicKey = base64.RawURLEncoding.EncodeToString(pub)
	privateKey = base64.RawURLEncoding.EncodeToString(priv)
	return
}
