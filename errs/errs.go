package errs

import (
	"errors"
)

var (
	WrongFile             = errors.New("wrong file")
	NetworkCheck          = errors.New("unsupported network")
	FileSaving            = errors.New("file saving failed")
	UpdateFsInfo          = errors.New("fs info update failed")
	WrongSignature        = errors.New("wrong signature")
	FileCheck             = errors.New("file check failed")
	ParseMultipartForm    = errors.New("parse multipart form failed")
	SpaceCheck            = errors.New("space check failed")
	NoSpace               = errors.New("not enough space")
	Internal              = errors.New("node internal error")
	InvalidArgument       = errors.New("invalid argument")
	StorageSystemNotFound = errors.New("storage filesystem not found")
)

// ====================================================================================
