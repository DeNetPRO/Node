package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"net"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"github.com/ethereum/go-ethereum/common"
)

var (
	NodeAddr   []byte
	PrivateKey []byte
)

// ====================================================================================

func EncryptAES(key, data []byte) ([]byte, error) {
	const logInfo = "shared.encryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil

}

// ====================================================================================

func DecryptAES(key, data []byte) ([]byte, error) {
	const logInfo = "shared.decryptAES->"
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}
	nonce, encrData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrData, err := gcm.Open(nil, nonce, encrData, nil)
	if err != nil {
		return nil, logger.CreateDetails(logInfo, err)
	}

	return decrData, nil
}

// ====================================================================================

func GetDeviceMacAddr() (string, error) {
	const logInfo = "shared.GetDeviceMacAddr->"
	var addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", logger.CreateDetails(logInfo, err)
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
		return nodeAddr, logger.CreateDetails(logInfo, err)
	}

	encrKey := sha256.Sum256([]byte(macAddr))

	encryptedAddr, err := EncryptAES(encrKey[:], addr.Bytes())
	if err != nil {
		return nodeAddr, logger.CreateDetails(logInfo, err)
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
		return nodeAddr, logger.CreateDetails(logInfo, err)
	}

	encrKey := sha256.Sum256([]byte(macAddr))

	accAddr, err := DecryptAES(encrKey[:], NodeAddr)
	if err != nil {
		return nodeAddr, logger.CreateDetails(logInfo, err)
	}

	return common.BytesToAddress(accAddr), nil
}

// ====================================================================================
