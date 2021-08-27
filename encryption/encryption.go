package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
)

var (
	PrivateKey []byte //encrypted private key
)

// ====================================================================================

//EncryptAES encrypts data using a provided key.
func EncryptAES(key, data []byte) ([]byte, error) {
	const location = "shared.encryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// ====================================================================================

//DecryptAES decrypts data using a provided key.
func DecryptAES(key, data []byte) ([]byte, error) {
	const location = "shared.decryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}
	nonce, encrData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrData, err := gcm.Open(nil, nonce, encrData, nil)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	return decrData, nil
}

// ====================================================================================

//GetDeviceMacAddr returns device's MAC address.
func GetDeviceMacAddr() (string, error) {
	const location = "shared.GetDeviceMacAddr->"
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", logger.CreateDetails(location, err)
	}

	for _, i := range interfaces {
		if !bytes.Equal(i.HardwareAddr, nil) {
			addr = i.HardwareAddr.String()
			break
		}
	}

	return addr, nil
}

// ====================================================================================
