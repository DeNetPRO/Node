package update

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/encryption"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	nodeAbi "git.denetwork.xyz/dfile/dfile-secondary-node/node_abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

type UpdatedFsInfo struct {
	NewFs            []string
	Nonce            string
	SignedFsRootHash string
}

func FsInfo(senderNodeAddr, storageAddr, signedFsRootHash, nonce string, fsHashes []string, nonce32, fsRootNonceBytes []byte) {

	const logInfo = "update.FsInfo->"

	concatFsHashes := ""

	for _, hash := range fsHashes {
		concatFsHashes += hash
	}

	hashesNonceBytes := append([]byte(concatFsHashes), nonce32...)

	hashesNonceSha := sha256.Sum256(hashesNonceBytes)

	encrKey := sha256.Sum256([]byte(senderNodeAddr))

	decryptedData, err := encryption.DecryptAES(encrKey[:], encryption.PrivateKey)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	accPrivKey, err := crypto.HexToECDSA(hex.EncodeToString(decryptedData))
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	signedFSTree, err := crypto.Sign(hashesNonceSha[:], accPrivKey)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	updatedFs := UpdatedFsInfo{
		NewFs:            fsHashes,
		Nonce:            nonce,
		SignedFsRootHash: signedFsRootHash,
	}

	updatedFsJson, err := json.Marshal(updatedFs)
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	nodeNft, err := blockchainprovider.GetNodeNFT()
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	totalNodes, err := nodeNft.TotalSupply(&bind.CallOpts{})
	if err != nil {
		logger.Log(logger.CreateDetails(logInfo, err))
	}

	for i := 0; i < int(totalNodes.Int64()); i++ {

		node, err := nodeNft.GetNodeById(&bind.CallOpts{}, big.NewInt(int64(i)))
		if err != nil {
			logger.CreateDetails(logInfo, err)
		}

		stringIP := getNodeIP(node)

		if stringIP == account.NodeIpAddr {
			continue
		}

		url := "http://" + stringIP + fmt.Sprint("/update_fs/", storageAddr, "/", senderNodeAddr, "/", hex.EncodeToString(signedFSTree))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(updatedFsJson))
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		go func(req *http.Request, stringIP string) {
			resp, _ := client.Do(req)

			if resp != nil {
				defer resp.Body.Close()

				if resp.Status != "200 OK" {
					logger.Log(logger.CreateDetails(logInfo, errors.New(stringIP+" fs wasn't updated")))
				}
			}
		}(req, stringIP)

	}

}

// ====================================================================================

func getNodeIP(nodeInfo nodeAbi.SimpleMetaDataDeNetNode) string {
	ipBuilder := strings.Builder{}
	for i, v := range nodeInfo.IpAddress {
		stringPart := strconv.Itoa(int(v))
		ipBuilder.WriteString(stringPart)

		if i < 3 {
			ipBuilder.WriteString(".")
		}
	}

	stringPort := strconv.Itoa(int(nodeInfo.Port))
	ipBuilder.WriteString(":")
	ipBuilder.WriteString(stringPort)

	return ipBuilder.String()
}
