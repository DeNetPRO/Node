package encryption_test

import (
	"testing"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"github.com/stretchr/testify/require"
)

var (
	testEncrKey   = []byte{25, 74, 62, 237, 207, 57, 204, 1, 136, 227, 96, 64, 103, 135, 205, 246, 169, 192, 122, 130, 197, 207, 77, 160, 123, 41, 16, 29, 126, 173, 58, 1}
	secretMessage = []byte{115, 101, 99, 114, 101, 116, 32, 109, 101, 115, 115, 97, 103, 101}
)

func TestEncryptDecrypt(t *testing.T) {

	encrOutput, err := encryption.EncryptAES(testEncrKey, secretMessage)
	if err != nil {
		t.Fatal(err)
	}

	decrOutput, err := encryption.DecryptAES(testEncrKey, encrOutput)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, secretMessage, decrOutput)

}
