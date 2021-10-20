package errs

import (
	"errors"
	"strings"
)

var (
	WrongFile             = errors.New("wrong file")
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

//Сromplatform error checking for file stat
func CheckStatErr(statErr error) error {
	if statErr == nil {
		return nil
	}

	errParts := strings.Split(statErr.Error(), ":")

	if len(errParts) == 3 && strings.Trim(errParts[2], " ") == "The system cannot find the file specified." {
		return nil
	}

	if len(errParts) == 2 && strings.Trim(errParts[1], " ") == "no such file or directory" {
		return nil
	}

	return statErr
}
