package errs

import (
	"errors"

	nodeTypes "github.com/DeNetPRO/src/node_types"
)

var errorList = nodeTypes.ErrList{
	FileName:      errors.New("wrong file"),
	Network:       errors.New("unsupported network"),
	FileSave:      errors.New("file saving failed"),
	FsUpdate:      errors.New("fs info update failed"),
	Signature:     errors.New("wrong signature"),
	FileCheck:     errors.New("file check failed"),
	Multipart:     errors.New("parse multipart form failed"),
	SpaceCheck:    errors.New("space check failed"),
	Space:         errors.New("not enough space"),
	Internal:      errors.New("node internal error"),
	Argument:      errors.New("invalid argument"),
	StorageSystem: errors.New("storage filesystem not found"),
}

// ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

func List() nodeTypes.ErrList {
	return errorList
}
