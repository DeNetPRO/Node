package dnetsignature

import "github.com/ethereum/go-ethereum/crypto"

// CheckDataSign checks if signature belongs to the sender.
func Check(spAddress string, signature []byte, hash [32]byte) error {
	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return err
	}

	signatureAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if spAddress != signatureAddress.String() {
		return err
	}

	return nil
}
