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
	PrivateKey []byte
)

// ====================================================================================

func EncryptAES(key, data []byte) ([]byte, error) {
	const logLoc = "shared.encryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil

}

// ====================================================================================

func DecryptAES(key, data []byte) ([]byte, error) {
	const logLoc = "shared.decryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}
	nonce, encrData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrData, err := gcm.Open(nil, nonce, encrData, nil)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	return decrData, nil
}

// ====================================================================================

func GetDeviceMacAddr() (string, error) {
	const logLoc = "shared.GetDeviceMacAddr->"
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", logger.CreateDetails(logLoc, err)
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
