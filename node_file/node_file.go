package nodefile

import (
	"encoding/json"
	"io"
	"os"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
)

//Read file by certain path
func Read(path string) (*os.File, []byte, error) {
	const location = "shared.ReadFile->"
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
	const location = "shared.ReadFromConsole->"

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
