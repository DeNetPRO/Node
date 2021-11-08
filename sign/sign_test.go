package sign_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/shared"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/sign"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCheckSignature(t *testing.T) {

	secrKeyHash := sha256.Sum256(shared.TestSecretKey)

	privateKeyBytes, err := encryption.DecryptAES(secrKeyHash[:], shared.TestPKHash)
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

	err = sign.Check(shared.TestAccAddr, hex.EncodeToString(signedData), hashData)
	if err != nil {
		t.Fatal(err)
	}
}
