package sign

import (
	"encoding/hex"

	"github.com/DeNetPRO/src/errs"
	"github.com/ethereum/go-ethereum/crypto"
)

// CheckDataSign checks if signature belongs to the sender.
func Check(signerAddress, signedData string, unsignedDataHash [32]byte) error {
	signature, err := hex.DecodeString(signedData)
	if err != nil {
		return err
	}

	sigPublicKey, err := crypto.SigToPub(unsignedDataHash[:], signature)
	if err != nil {
		return err
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if signerAddress != signatureAddress.String() {
		return errs.List().Signature
	}

	return nil
}
