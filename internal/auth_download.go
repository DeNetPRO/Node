package internal

import (
	"errors"
	"math/rand"
	"time"

	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"github.com/valyala/fasthttp"
)

func AuthDownload(spAddress, fileKey string) (*NodeAddressResponse, error) {
	const logLoc = "internal.AuthDownload->"

	nodes, err := shared.GetFileInfoNodes(spAddress, fileKey)
	if err != nil {
		return nil, logger.CreateDetails(logLoc, err)
	}

	fastReq := fasthttp.AcquireRequest()
	fastResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(fastReq)
	defer fasthttp.ReleaseResponse(fastResp)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nodes), func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })

	for _, node := range nodes {
		fastReq.Reset()
		fastResp.Reset()

		nodeAddress, err := checkFileOnNode(node, spAddress, fileKey, fastReq, fastResp)
		if err != nil {
			err = shared.RemoveConnectionNode(node)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
			}

			err = shared.RemoveFileKeySavedNode(spAddress, fileKey, node)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
			}

			continue
		}

		return nodeAddress, nil
	}

	return nil, logger.CreateDetails(logLoc, errors.New("no file"))
}
