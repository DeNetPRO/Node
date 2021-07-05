package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"dfile-secondary-node/logger"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/ethereum/go-ethereum/common"
)

var (
	NodeAddr []byte
)

// ====================================================================================

func encryptAES(key, data []byte) ([]byte, error) {
	const logInfo = "shared.encryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil

}

// ====================================================================================

func decryptAES(key, data []byte) ([]byte, error) {
	const logInfo = "shared.decryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}
	nonce, encrData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrData, err := gcm.Open(nil, nonce, encrData, nil)
	if err != nil {
		return nil, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	return decrData, nil
}

// ====================================================================================

func GetDeviceMacAddr() (string, error) {
	const logInfo = "shared.GetDeviceMacAddr->"
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
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

func EncryptNodeAddr(addr common.Address) ([]byte, error) {
	const logInfo = "shared.EncryptNodeAddr->"
	var nodeAddr []byte

	macAddr, err := GetDeviceMacAddr()
	if err != nil {
		return nodeAddr, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	encrKey := sha256.Sum256([]byte(macAddr))

	encryptedAddr, err := encryptAES(encrKey[:], addr.Bytes())
	if err != nil {
		return nodeAddr, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	return encryptedAddr, nil
}

// ====================================================================================

func DecryptNodeAddr() (common.Address, error) {
	const logInfo = "shared.DecryptNodeAddr->"
	var nodeAddr common.Address

	if len(NodeAddr) == 0 {
		return nodeAddr, errors.New("empty address")
	}

	macAddr, err := GetDeviceMacAddr()
	if err != nil {
		return nodeAddr, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	encrKey := sha256.Sum256([]byte(macAddr))

	accAddr, err := decryptAES(encrKey[:], NodeAddr)
	if err != nil {
		return nodeAddr, fmt.Errorf("%s %w", logInfo, logger.GetDetailedError(err))
	}

	return common.BytesToAddress(accAddr), nil
}

// ====================================================================================
