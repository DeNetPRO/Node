package rpcserver

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"sync"

	"github.com/DeNetPRO/src/config"
	"github.com/DeNetPRO/src/hash"
	"github.com/DeNetPRO/src/logger"
	"github.com/DeNetPRO/src/networks"
	"github.com/DeNetPRO/src/paths"
	"github.com/DeNetPRO/src/pb"
	"github.com/DeNetPRO/src/sign"
	spFiles "github.com/DeNetPRO/src/sp_files"

	fsysInfo "github.com/DeNetPRO/src/fsys_info"

	nodeFile "github.com/DeNetPRO/src/node_file"

	nodeTypes "github.com/DeNetPRO/src/node_types"

	tstpkg "github.com/DeNetPRO/src/tst_pkg"

	"google.golang.org/grpc"
)

type rpcServer struct {
	pb.UnimplementedNodeServiceServer
}

var mutex sync.Mutex

func Start(port string) error {

	const location = "rpcserver.Start ->"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	s := grpc.NewServer()

	pb.RegisterNodeServiceServer(s, &rpcServer{})

	fmt.Println("starting rpc server on port", port)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	s.GracefulStop()

	if tstpkg.Data().TestMode {
		os.RemoveAll(paths.List().WorkDir)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (r *rpcServer) UploadFile(stream pb.NodeService_UploadFileServer) error {

	const location = "rpcserver.UploadFile ->"

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	err = sign.Check(req.SpAddress, req.SignedAddress, sha256.Sum256([]byte(req.SpAddress)))
	if err != nil {
		return err
	}

	err = networks.Check(req.Network)
	if err != nil {
		return errors.New("unsupported network")
	}

	pathToSpFiles := filepath.Join(paths.List().Storages[0], req.Network, req.SpAddress)

	dirStat, err := os.Stat(pathToSpFiles)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Log(logger.MarkLocation(location, err))
		return errors.New("couldn't check SP files")
	}

	err = checkAndReserveSpace(req.FileSize)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return errors.New("couldn't reserve space")
	}

	if dirStat == nil {
		err = os.MkdirAll(pathToSpFiles, 0700)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			return errors.New("couldn't create storage")
		}
	}

	for {

		req, err = stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		err = spFiles.SaveChunk(pathToSpFiles, req.FileName, req.ChunkData)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			return errors.New("couldn't save file")
		}

		fmt.Println("saved file:", req.FileName)

	}

	stream.SendAndClose(&pb.Response{Msg: "saved"})

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (r *rpcServer) DownloadFile(req *pb.DownloadRequest, srv pb.NodeService_DownloadFileServer) error {

	err := sign.Check(req.SpAddress, req.SignedAddress, sha256.Sum256([]byte(req.SpAddress)))
	if err != nil {
		return err
	}

	err = networks.Check(req.Network)
	if err != nil {
		return err
	}

	pathToStorage := filepath.Join(paths.List().Storages[0], req.Network, req.SpAddress)

	for _, fileName := range req.FileNames {

		pathToFile := filepath.Join(pathToStorage, fileName)

		_, err := os.Stat(pathToFile)
		if err != nil {
			return err
		}

		bytes, err := os.ReadFile(pathToFile)
		if err != nil {
			return err
		}

		err = srv.Send(&pb.DownloadResponse{ChunkData: bytes})
		if err != nil {
			return err
		}

		fmt.Println("serving file:", fileName)

	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (r *rpcServer) GatewayDownloadFile(req *pb.GatewayDownloadRequest, srv pb.NodeService_GatewayDownloadFileServer) error {

	err := sign.Check(req.GatewayAddress, req.SignedGatewayAddress, sha256.Sum256([]byte(req.GatewayAddress)))
	if err != nil {
		return err
	}

	pathToStorage := filepath.Join(paths.List().Storages[0], req.Network, req.SpAddress)

	for _, fileName := range req.FileNames {

		pathToFile := filepath.Join(pathToStorage, fileName)

		_, err := os.Stat(pathToFile)
		if err != nil {
			return err
		}

		bytes, err := os.ReadFile(pathToFile)
		if err != nil {
			return err
		}

		err = srv.Send(&pb.DownloadResponse{ChunkData: bytes})
		if err != nil {
			return err
		}

		fmt.Println("serving file:", fileName)

	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func (r *rpcServer) UpdateFs(ctx context.Context, req *pb.FsInfo) (*pb.Response, error) {

	const location = "rpcserver.UpdateFs ->"

	sort.Strings(req.NewFs)

	fsRootHash, fsTree, err := hash.CalcRoot(req.NewFs)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return &pb.Response{Msg: "failed"}, errors.New("couldn't save fsys info")
	}

	fsRootBytes, err := hex.DecodeString(fsRootHash)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return &pb.Response{Msg: "failed"}, errors.New("couldn't save fsys info")
	}

	nonceBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(nonceBytes, req.Nonce)

	alignBytes := make([]byte, 28) // need to align nonce and storage info to 32 bytes for smart contract

	alignedNonceBytes := append(alignBytes, nonceBytes...)

	storageBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(storageBytes, req.Storage)

	alignedStorage := append(alignBytes, storageBytes...)

	fsRootStorageBytes := append(fsRootBytes, alignedStorage...)

	fsRootStorageNonceBytes := append(fsRootStorageBytes, alignedNonceBytes...)

	err = sign.Check(req.SpAddress, req.Signature, sha256.Sum256(fsRootStorageNonceBytes))
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return &pb.Response{Msg: "failed"}, errors.New("couldn't save fsys info")
	}

	err = fsysInfo.Save(req, fsTree)
	if err != nil {
		logger.Log(logger.MarkLocation(location, err))
		return &pb.Response{Msg: "failed"}, errors.New("couldn't save fsys info")
	}

	return &pb.Response{Msg: "updated"}, nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

//Checks if space enough for uploading a file and reserves space

func checkAndReserveSpace(fileSize uint32) error {
	const location = "rpcserver.checkAndReserveSpace"

	var nodeConfig nodeTypes.Config

	mutex.Lock()
	defer mutex.Unlock()

	confFile, fileBytes, err := nodeFile.Read(paths.List().ConfigFile)
	if err != nil {
		mutex.Unlock()
		return logger.MarkLocation(location, err)
	}
	defer confFile.Close()

	err = json.Unmarshal(fileBytes, &nodeConfig)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	totalNodeSpace := int64(nodeConfig.StorageLimit * 1024 * 1024 * 1024) // convert to bytes

	nodeConfig.UsedStorageSpace += int64(fileSize)

	if nodeConfig.UsedStorageSpace > totalNodeSpace {
		return logger.MarkLocation(location, err)
	}

	avaliableSpaceLeft := totalNodeSpace - nodeConfig.UsedStorageSpace

	if avaliableSpaceLeft < int64(1024*1024*100) { // 100 Mib
		fmt.Println("Shared storage memory is running low,", avaliableSpaceLeft/(1024*1024), "MB of space is avaliable")
		fmt.Println("You may need additional space for storing data. Total shared space can be changed in account configuration")
	}

	err = config.Save(confFile, nodeConfig)
	if err != nil {
		return logger.MarkLocation(location, err)
	}

	return nil
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
