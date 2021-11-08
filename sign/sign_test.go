package sign_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	testPasswd  = "testPasswd"
	testAccAddr string
)

func TestMain(m *testing.M) {
	shared.TestModeOn()
	defer shared.TestModeOff()

	err := paths.Init()
	if err != nil {
		log.Fatal(err)
	}

	accAddr, _, err := account.Create(testPasswd)
	if err != nil {
		log.Fatal(err)
	}

	testAccAddr = accAddr

	exitVal := m.Run()

	err = os.RemoveAll(paths.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}

func TestCheckSignature(t *testing.T) {

	secrKeyHash := sha256.Sum256(encryption.SecretKey)

	privateKeyBytes, err := encryption.DecryptAES(secrKeyHash[:], encryption.EncryptedPK)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 100)
	rand.Read(data)

	hashData := sha256.Sum256(data)

	signedData, err := crypto.Sign(hashData[:], privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = sign.Check(testAccAddr, hex.EncodeToString(signedData), hashData)
	if err != nil {
		t.Fatal(err)
	}
}
