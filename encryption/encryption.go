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
	const actLoc = "shared.encryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil

}

// ====================================================================================

func DecryptAES(key, data []byte) ([]byte, error) {
	const actLoc = "shared.decryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}
	nonce, encrData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrData, err := gcm.Open(nil, nonce, encrData, nil)
	if err != nil {
		return nil, logger.CreateDetails(actLoc, err)
	}

	return decrData, nil
}

// ====================================================================================

func GetDeviceMacAddr() (string, error) {
	const actLoc = "shared.GetDeviceMacAddr->"
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", logger.CreateDetails(actLoc, err)
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
