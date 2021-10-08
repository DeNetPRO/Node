package test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"github.com/stretchr/testify/require"
)

type FileSendInfo struct {
	Hash string `json:"hash"`
	Body []byte `json:"body"`
}

var (
	accountPassword           = "123"
	accountAddress            string
	nodeAddress               []byte
	ErrorInvalidPassword      = errors.New(" could not decrypt key with given password")
	configPath                string
	storagePath               string
	testFileName              = "file"
	fileSize                  int64
	testFilePath              string
	NodeTestPrivateKey        = "16f98d96422dd7f21965755bd64c9dcd9cfc5d36e029002d9cc579f42511c7ed"
	storageProviderPrivateKey = "0a9fb845e346f74227d2ddf0b85dedb4ccddee33e9b8d0f6f4828a7a2dcf9509"
	storageProviderAddress    = "0x3429cC113ABf4DEc8ECA64A713761F90A000dDfB"
)

func TestMain(m *testing.M) {
	err := os.Setenv("DENET_TEST", "1")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Unsetenv("DENET_TEST")

	paths.WorkDirName = "denet-test"

	err = paths.Init()
	if err != nil {
		log.Fatal("Fatal Error: couldn't locate home directory")
	}
	exitVal := m.Run()

	err = os.RemoveAll(paths.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestEmptyAccountListBeforeCreating(t *testing.T) {

	accs := account.List()
	want := 0
	get := len(accs)

	require.Equal(t, want, get)
}

func TestAccCreate(t *testing.T) {
	_, _, err := account.Create(accountPassword)
	if err != nil {
		t.Error(err)
	}

	stat, err := os.Stat(paths.AccsDirPath)
	if err != nil {
		t.Error(err)
	}

	accs := account.List()
	if len(accs) != 1 {
		t.Error("Wrong accs count, must be one", accs)
	}

	fmt.Println(stat)
}

// func TestLoginAccountWithCorrectAddressAndPassword(t *testing.T) {
// 	account, err := account.Login(accountAddress, accountPassword)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	require.Equal(t, accountAddress, account.Address.String())
// }

// func TestLoginAccountWithInvalidPassword(t *testing.T) {
// 	_, err := account.Login(accountAddress, "invalid")
// 	want := ErrorInvalidPassword

// 	splitErr := strings.Split(err.Error(), "->")

// 	require.EqualError(t, want, splitErr[len(splitErr)-1])
// }

// func TestLoginAccountWithUnknownAddress(t *testing.T) {
// 	unknownAddress := "accountAddress"
// 	_, err := account.Login(unknownAddress, accountPassword)
// 	want := errors.New(" accountAddress address is not found")
// 	splitErr := strings.Split(err.Error(), "->")

// 	require.EqualError(t, want, splitErr[len(splitErr)-1])
// }

// func TestCheckRightPassword(t *testing.T) {
// 	err := account.CheckPassword(accountPassword, accountAddress)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestImportAccount(t *testing.T) {
// 	accountAddress, c, err := account.Import()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if accountAddress == "" {
// 		t.Errorf("import account address must not to b	e empty")
// 	}

// 	wantConfig := config.NodeConfig{
// 		Address: accountAddress,
// 		Network: blckChain.Network,

// 		AgreeSendLogs: true,
// 	}

// 	require.Equal(t, wantConfig, c)
// }

// func TestCheckSignature(t *testing.T) {
// 	macAddress, err := encryption.GetDeviceMacAddr()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	encrForKey := sha256.Sum256([]byte(macAddress))
// 	privateKeyBytes, err := encryption.DecryptAES(encrForKey[:], encryption.PrivateKey)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	privateKey, err := crypto.ToECDSA(privateKeyBytes)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	data := make([]byte, 100)
// 	rand.Seed(time.Now().Unix())
// 	rand.Read(data)

// 	hashData := sha256.Sum256(data)

// 	signedData, err := crypto.Sign(hashData[:], privateKey)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = sign.Check(accountAddress, signedData, hashData)
// 	if err != nil {
// 		t.Error(encrForKey)
// 	}
// }

// func TestRestoreNodeMemory(t *testing.T) {
// 	fileSize := 1024 * 1024

// 	confFile, nodeConfig, err := getConfig()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	want := nodeConfig.UsedStorageSpace

// 	nodeConfig.UsedStorageSpace += int64(fileSize)

// 	err = config.Save(confFile, *nodeConfig)
// 	if err != nil {
// 		confFile.Close()
// 		t.Error(err)
// 	}

// 	confFile.Close()

// 	meminfo.Restore(configPath, fileSize)

// 	confFile, nodeConfig, err = getConfig()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	confFile.Close()

// 	require.Equal(t, want, nodeConfig.UsedStorageSpace)
// }

// func TestUpload(t *testing.T) {
// 	createFilesForTest()

// 	file, err := os.Open(testFilePath)
// 	if err != nil {
// 		return
// 	}

// 	defer file.Close()

// 	const oneMB = 1048548 // actually less than 1MB because we need additional 12 bytes for nonce and 16 bytes is max overhead for encoding
// 	const eightKB = 8192

// 	fileChunk := make([]byte, oneMB-7) // reserve 7 bytes for added zeros info
// 	lenFileChunk := len(fileChunk)
// 	count := 0
// 	oneMBHashes := make([]string, 0, fileSize/oneMB+1)
// 	eightKBHashes := make([]string, 0, 128)
// 	partToEncode := make([]byte, 0, oneMB)
// 	addedZeros := make([]byte, 7)

// 	fileSendInfo := make(map[string]*FileSendInfo)

// 	for {
// 		bytes, err := file.Read(fileChunk)
// 		if err != nil {
// 			if err != io.EOF {
// 				return
// 			}

// 			break
// 		}

// 		if bytes < lenFileChunk {
// 			missPart := make([]byte, lenFileChunk-bytes)
// 			fileChunk = append(fileChunk[:bytes], missPart...)

// 			misspartLenInBytes := []byte(fmt.Sprint(len(missPart)))

// 			for i := 0; i < len(misspartLenInBytes); i++ {
// 				addedZeros[i] = misspartLenInBytes[i]
// 			}
// 		}

// 		partToEncode = append(partToEncode, addedZeros...)
// 		partToEncode = append(partToEncode, fileChunk...)

// 		key := CreatePartEncrKey("dir", fmt.Sprint("part_", count))

// 		encryptedPart, err := encryption.EncryptAES(key, partToEncode)
// 		if err != nil {
// 			return
// 		}

// 		partToEncode = partToEncode[:0]

// 		for i := 0; i < len(encryptedPart); i += eightKB {
// 			hSum := sha256.Sum256(encryptedPart[i : i+eightKB])
// 			eightKBHashes = append(eightKBHashes, hex.EncodeToString(hSum[:]))
// 		}

// 		oneMBHash, _, err := hash.CalcRoot(eightKBHashes)
// 		if err != nil {
// 			return
// 		}

// 		eightKBHashes = eightKBHashes[:0]
// 		oneMBHashes = append(oneMBHashes, oneMBHash)

// 		fileSendInfo[oneMBHash] = &FileSendInfo{
// 			Hash: oneMBHash,
// 			Body: encryptedPart,
// 		}

// 		count++
// 	}

// 	pipeConns := fasthttputil.NewPipeConns()
// 	pr := pipeConns.Conn1()
// 	pw := pipeConns.Conn2()

// 	writer := multipart.NewWriter(pw)

// 	go prepareFileBeforeUpload(writer, pw, oneMBHashes, fileSendInfo)

// 	endpoint := "/upload/" + strconv.Itoa(int(fileSize))

// 	req, err := http.NewRequest("POST", endpoint, pr)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	rr := httptest.NewRecorder()

// 	router := mux.NewRouter()
// 	router.HandleFunc("/upload/{size}", server.SaveFiles)
// 	router.ServeHTTP(rr, req)

// 	require.Equal(t, "OK", rr.Body.String())

// 	for _, fileName := range oneMBHashes {
// 		path := filepath.Join(paths.AccsDirPath, shared.NodeAddr.String(), paths.StorageDirName, storageProviderAddress, fileName)
// 		_, err := os.Stat(path)
// 		if err != nil {
// 			t.Errorf("%v not saved", fileName)
// 		}
// 	}
// }

// func getConfig() (*os.File, *config.NodeConfig, error) {
// 	confFile, fileBytes, err := nodefile.Read(configPath)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var nodeConfig *config.NodeConfig

// 	err = json.Unmarshal(fileBytes, &nodeConfig)
// 	if err != nil {
// 		confFile.Close()
// 		return nil, nil, err
// 	}

// 	return confFile, nodeConfig, nil
// }

// func createFilesForTest() {
// 	rand.Seed(time.Now().Unix())

// 	fileSize = 1024 * 1024 * 10

// 	b := make([]byte, fileSize)

// 	testFilePath = filepath.Join(paths.AccsDirPath, testFileName)

// 	f, err := os.Create(testFilePath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	rand.Read(b)

// 	f.Write(b)
// 	f.Close()
// }

// func CreatePartEncrKey(dirName, part string) []byte {
// 	h := sha256.Sum256([]byte(dirName))

// 	strKey := hex.EncodeToString(h[:])

// 	encrKey := sha256.Sum256([]byte(fmt.Sprint(strKey, part)))

// 	return encrKey[:]
// }

// func prepareFileBeforeUpload(writer *multipart.Writer, pw net.Conn, oneMBHashes []string, fileSendInfo map[string]*FileSendInfo) {
// 	defer pw.Close()

// 	var fileRootHash string
// 	var err error

// 	if len(oneMBHashes) == 1 {
// 		fileRootHash = oneMBHashes[0]
// 	} else {
// 		sort.Strings(oneMBHashes)
// 		fileRootHash, _, err = hash.CalcRoot(oneMBHashes)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	oneMBHashes = append(oneMBHashes, fileRootHash)

// 	var wholeRootHash string
// 	sort.Strings(oneMBHashes)
// 	wholeRootHash, _, err = hash.CalcRoot(oneMBHashes)
// 	if err != nil {
// 		return
// 	}

// 	nonceInt := time.Now().Unix()

// 	err = writer.WriteField("address", storageProviderAddress)
// 	if err != nil {
// 		return
// 	}

// 	err = writer.WriteField("nonce", fmt.Sprint(nonceInt))
// 	if err != nil {
// 		return
// 	}

// 	nonceHex := strconv.FormatInt(nonceInt, 16)

// 	nonceBytes, err := hex.DecodeString(nonceHex)
// 	if err != nil {
// 		return
// 	}

// 	nonce32 := make([]byte, 32-len(nonceBytes))
// 	nonce32 = append(nonce32, nonceBytes...)

// 	fsRootBytes, err := hex.DecodeString(wholeRootHash)
// 	if err != nil {
// 		return
// 	}
// 	fsRootNonceBytes := append(fsRootBytes, nonce32...)

// 	hash := sha256.Sum256(fsRootNonceBytes)

// 	pk, _ := crypto.HexToECDSA(storageProviderPrivateKey)

// 	signedFSRootHash, err := crypto.Sign(hash[:], pk)
// 	if err != nil {
// 		return
// 	}

// 	err = writer.WriteField("fsRootHash", hex.EncodeToString(signedFSRootHash))
// 	if err != nil {
// 		return
// 	}

// 	for _, hash := range oneMBHashes {
// 		err := writer.WriteField("fs", hash)
// 		if err != nil {
// 			return
// 		}

// 		if hash != fileRootHash {
// 			filePart, err := writer.CreateFormFile("files", hash)
// 			if err != nil {
// 				return
// 			}

// 			filePart.Write(fileSendInfo[hash].Body)
// 		}
// 	}

// 	writer.Close()
// }
