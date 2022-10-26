package sign_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCheckSignature(t *testing.T) {

	privateKeyBytes, err := encryption.DecryptAES(tstpkg.Data().EncrKey, tstpkg.Data().PKHash)
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

	err = sign.Check(tstpkg.Data().AccAddr, hex.EncodeToString(signedData), hashData)
	if err != nil {
		t.Fatal(err)
	}
}
