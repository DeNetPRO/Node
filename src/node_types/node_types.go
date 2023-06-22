package types

type Config struct {
	Address              string            `json:"nodeAddress"`
	IpAddress            string            `json:"ipAddress"`
	HTTPPort             string            `json:"portHTTP"`
	Network              string            `json:"network"`
	RPC                  map[string]string `json:"rpc"`
	StorageLimit         int               `json:"storageLimit"`
	StoragePaths         []string          `json:"storagePaths"`
	UsedStorageSpace     int64             `json:"usedStorageSpace"`
	SendBugReports       bool              `json:"sendBugReports"`
	RegisteredInNetworks map[string]bool   `json:"registeredInNetworks"`
}

type NtwrkParams struct {
	TRX  string
	RPC  string
	NODE string
	PoS  string
	ERC  string
}

type UpdatedFsInfo struct {
	NewFs                 []string `json:"newFs"`
	Nonce                 int64    `json:"nonce"`
	SignedFsRootNonceHash string   `json:"signedFsRootNonceHash"`
}

type StatsInfoData struct {
	Type       string
	FileSize   int64
	RemoteAddr string
	Network    string
}

type ReqData struct {
	RequesterAddr string
	FileName      string
	FsTreeHash    string
}

type StorageProviderData struct {
	Nonce        uint32     `json:"nonce"`
	Storage      uint32     `json:"storage"`
	SignedFsInfo string     `json:"signedFsRoot"`
	Tree         [][][]byte `json:"tree"`
}

type NodesResponse struct {
	Nodes []string `json:"nodes"`
}

type FileSendInfo struct {
	Hash string `json:"hash"`
	Body []byte `json:"body"`
}

type ErrList struct {
	FileName      error
	Network       error
	FileSave      error
	FsUpdate      error
	Signature     error
	FileCheck     error
	Multipart     error
	SpaceCheck    error
	Space         error
	Internal      error
	Argument      error
	StorageSystem error
}

type Paths struct {
	WorkDir      string
	AccsDir      string
	ConfigDir    string
	ConfigFile   string
	UpdateDir    string
	SysDir       string
	SpFsFilename string
	Storages     []string
}
