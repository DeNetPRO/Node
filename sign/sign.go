package sign

import (
	"encoding/hex"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"
	"github.com/ethereum/go-ethereum/crypto"
)

// CheckDataSign checks if signature belongs to the sender.
func Check(spAddress, signedData string, unsignedDataHash [32]byte) error {
	signature, err := hex.DecodeString(signedData)
	if err != nil {
		return err
	}

	sigPublicKey, err := crypto.SigToPub(unsignedDataHash[:], signature)
	if err != nil {
		return err
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return errs.WrongSignature
	}

	return nil
}
