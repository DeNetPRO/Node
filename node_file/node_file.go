package nodefile

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
)

//Read file by certain path
func Read(path string) (*os.File, []byte, error) {
	const location = "nodefile.ReadFile->"
	f, err := os.OpenFile(path, os.O_RDWR, 0700)
	if err != nil {
		return nil, nil, logger.CreateDetails(location, err)
	}

	fBytes, err := io.ReadAll(f)
	if err != nil {
		f.Close()
		return nil, nil, logger.CreateDetails(location, err)
	}

	return f, fBytes, nil
}

// ====================================================================================

func Write(f *os.File, data interface{}) error {
	const location = "nodefile.ReadFromConsole->"

	js, err := json.Marshal(data)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	err = f.Truncate(0)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	_, err = f.Write(js)
	if err != nil {
		return logger.CreateDetails(location, err)
	}

	f.Sync()

	return nil
}

// ====================================================================================

func ReadDirFiles(path string) ([]fs.FileInfo, error) {

	const location = "nodefile.ReadDirFiles->"

	dir, err := os.Open(path)
	if err != nil {
		return nil, logger.CreateDetails(location, err)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		logger.CreateDetails(location, err)
		return nil, logger.CreateDetails(location, err)
	}

	return files, nil
}
