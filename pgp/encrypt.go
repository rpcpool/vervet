package pgp

import (
	"bytes"
	"compress/gzip"
	_ "crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
)

func EncryptB64(entity *openpgp.Entity, message []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	encryptorWriter, err := openpgp.Encrypt(buf, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating entity for encryption: %v", err)
	}
	messageReader := bytes.NewReader(message)

	_, err = io.Copy(encryptorWriter, messageReader)
	if err != nil {
		return []byte{}, fmt.Errorf("Error writing data to encryptor: %v", err)
	}
	encryptorWriter.Close()

	// Encode the encrypted message to base64
	encodedMessage := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Return the base64 encoded message
	return []byte(encodedMessage), nil
}

func Encrypt(entity *openpgp.Entity, message []byte) ([]byte, error) {
	// Create buffer to write output to
	buf := new(bytes.Buffer)

	// Create encoder
	encoderWriter, err := armor.Encode(buf, "Message", make(map[string]string))
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating OpenPGP armor: %v", err)
	}

	// Create encryptor with encoder
	encryptorWriter, err := openpgp.Encrypt(encoderWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating entity for encryption: %v", err)
	}

	// Create compressor with encryptor
	compressorWriter, err := gzip.NewWriterLevel(encryptorWriter, gzip.BestCompression)
	if err != nil {
		return []byte{}, fmt.Errorf("Invalid compression level: %v", err)
	}

	// Write message to compressor
	messageReader := bytes.NewReader(message)
	_, err = io.Copy(compressorWriter, messageReader)
	if err != nil {
		return []byte{}, fmt.Errorf("Error writing data to compressor: %v", err)
	}

	compressorWriter.Close()
	encryptorWriter.Close()
	encoderWriter.Close()

	// Return buffer output - an encoded, encrypted, and compressed message
	return buf.Bytes(), nil
}
