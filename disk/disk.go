package disk

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
)

var fileSystemTypes = map[string]bool{
	"ext4": true,
	"ext3": true,
	"ext2": true,
}

const DeNetDisk = "DeNet-Disk.img"

func InitStorageCapacity(gb int) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fmt.Println(homeDir)

	mntPath := filepath.Join(homeDir, "mnt", "DeNet")
	err = os.MkdirAll(mntPath, 0700)
	if err != nil {
		return err
	}

	fmt.Println(mntPath)
	fsType, err := getFileSystemTypeInHomeDir(homeDir)
	if err != nil {
		return err
	}

	fmt.Println(fsType)
	diskPath := filepath.Join(mntPath, DeNetDisk)
	fmt.Println(diskPath)
	err = initDisk(diskPath, fsType, gb)
	if err != nil {
		return err
	}

	err = mountDisk(diskPath, paths.WorkDirPath)
	if err != nil {
		return err
	}

	return nil
}

func getFileSystemTypeInHomeDir(homeDir string) (string, error) {

	cmd := exec.Command("df", homeDir, "-T")
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	dfSplit := strings.Split(buffer.String(), "\n")
	params := strings.Split(dfSplit[1], " ")

	for _, v := range params {
		if string(v) == " " {
			continue
		}

		_, ok := fileSystemTypes[string(v)]
		if ok {
			return string(v), nil
		}
	}

	return "", fmt.Errorf("not found")
}

func initDisk(diskPath, fsType string, sizeGB int) error {
	err := createDisk(diskPath, sizeGB)
	if err != nil {
		return err
	}

	return createFileSystem(diskPath, fsType)
}

func createDisk(diskPath string, size int) error {
	file, err := os.Create(diskPath)
	if err != nil {
		return err
	}

	defer file.Close()

	return file.Truncate(int64(size * 1024 * 1024 * 1024))
}

func createFileSystem(src, fsType string) error {
	cmd := "mkfs." + fsType
	return exec.Command(cmd, src).Run()
}

func mountDisk(src, target string) error {
	return exec.Command("mount", src, target).Run()
}
